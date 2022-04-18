
class BillsGenerator {
    target = undefined;
    target_id = undefined;
    header = undefined;
    consumer_id = 0;
    collection = {};
    bills = [];
    buffer = {};
    section = {};
    draft = [];
    objects = [];

    constructor(target, header, buffer, section, accounts) {
        let self = this;
        self.header = header;
        self.target_id = target;
        self.collection = accounts;
        self.buffer = buffer;
        self.section = section;
        self.consumer_id = accounts.consume((action, account, id, href)=>{
            if (action === "add") {
                account.objects.forEach((object) => {
                    self.objects.push({
                        object_id: object.object_id,
                        object: object,
                        account: account,
                    });
                    self.addAccountObject(account, object);
                })
            }
        });
        self.updateTargetInfo();
    }

    updateTargetInfo() {
        let self = this;
        manager.GetTarget(self.target_id, (target) => {
            self.target = target;
            self.header.content({
                tag: "div",
                class: "row",
                content: [
                    {tag: "div", class: ["col", "m6", "s12"], content: target.comment + " [" + target.cost + "руб.]"},
                    {tag: "div", id: "bills-controls-btns", class: ["col", "m6", "s12"], content: [
                            {tag: "button", class: "btn", content: "Оформить", onClick: () => self.apply()},
                            {tag: "button", class: ["btn", "red"], content: "Завершить", onClick: () => self.final()}
                        ]},
                    {tag: "div", class: ["row"], content: {tag: "div", class: ["col", "m12"], content: [
                                {tag: "div", id: "progress", class: ["progress"], content: {tag: "div", class: ["determinate"], style: "width: 100%"}}
                            ]}}
                ]
            });
            if (!target.closed) {
                $("#bills-controls-btns").show(200);
            } else {
                $("#bills-controls-btns").hide(200);
                $("#bill-generator-buffer").hide(200);
                $("#bills-prepared-header a").hide(200);
                $("#bill-generator-section").attr("class", "com s12");
            }
        }, (message, status) => {
            toast(message);
        })
    }

    fillBills(bills){
        this.bills = bills;
        bills.forEach((bill) => this.appendBill(bill))
    }

    appendBill(bill){
        let self = this;
        if (self.objects.find((o) => o.object_id === bill.object_id)) {
            $("#" + bill.object_id).remove();
        }
        self.section.append({
            tag: "a",
            id: bill.bill_id,
            href: "#bill:bill_id=" + bill.bill_id,
            class: ["collection-item", "truncate", "black-text"],
            content: [
                {tag : "span", id: "head-for-" + bill.object_id, content: self.getHeadFor(bill.object_id)},
                {tag: "br"},
                {tag: "span", id: "bill-for-" + bill.object_id, content: (new Date(bill.formed).toLocaleDateString()) + " " + bill.bill + " руб."}
            ]
        });
        self.draft = self.draft.filter((d) => d.object.object_id !== bill.object_id);
    }

    prependBillDraft(account, object){
        let self = this;
        self.section.prepend({
            tag: "a",
            id: object.object_id,
            href: "#!",
            onClick: ()=>{
                self.fromBill(account, object);
                return false;
            },
            class: ["collection-item", "truncate", "draft-item"],
            content: getObjectShortAddress(object) + "<br />" + getFirstPersonName(account)
        });
    }

    getHeadFor(object_id) {
        let object = this.objects.find((o) => o.object_id === object_id);
        if (object) {
            return getObjectShortAddress(object.object) + " " + getFirstPersonName(object.account)
        }
        return ". . ."
    }

    addAccountObject(account, object){
        let self = this;
        if (this.bills.find((bill) => bill.object_id === object.object_id)) {
            let head = new Render("#head-for-" + object.object_id);
            head.content(self.getHeadFor(object.object_id));
            return
        }
        self.buffer.prepend({
            tag: "a",
            id: object.object_id,
            href: "#!",
            onClick: ()=>{
                self.toBills(account, object);
                return false;
            },
            class: ["collection-item", "truncate"],
            content: getObjectShortAddress(object) + "<br />" + getFirstPersonName(account)
        });
    }

    toBills(account, object) {
        $("#" + object.object_id).remove();
        this.prependBillDraft(account, object);
        this.draft.push({
            account: account,
            object: object
        });
    }

    draftAll() {
        let self = this;
        $("#bills-accounts-list .collection-item").trigger('click');
    }

    undraftAll() {
        let self = this;
        $("#bills-prepared-list .draft-item").trigger('click');
    }

    fromBill(account, object) {
        $("#" + object.object_id).remove();
        this.addAccountObject(account, object);
        this.draft = this.draft.filter((d) => d.object.object_id !== object.object_id);
    }

    apply() {
        let self = this;
        if (self.target.closed) {
            return
        }
        let progress = new Render("#progress");
        let allAllCount = 0;
        let doneCount = 0;
        let buckets = this.draft.reduce((c, el) => {
            allAllCount++;
            let l = c.length;
            let buck = c[l-1];
            if (buck.length > 9) {
                c.push([]);
                buck = c[l];
                l++;
            }
            buck.push(el);
            c[l-1] = buck;
            return c;
        }, [[]]);
        let cur = 0;
        let process = () => {};
        let continuePlease = () => {
            process();
            cur++;
        };
        process = () => {
            if (buckets.length <= cur) {
                $("#bills-controls-btns").show(200);
                return
            }
            let allCnt = buckets[cur].length;
            buckets[cur].forEach((draft) => {
                manager.CreateBill(
                    draft.account.account_id,
                    {
                        formed: new Date(),
                        object_id: draft.object.object_id,
                        year: self.target.period.year,
                        month: self.target.period.month,
                        target_id: self.target.target_id,
                        type: self.target.type,
                        bill: self.target.cost,
                    },
                    (bill) => {
                        allCnt--;
                        doneCount++;
                        let perc = (doneCount / allAllCount) * 100;
                        progress.content(
                            {tag: "div", class: ["determinate"], style: "width: " + perc + "%"}
                        );
                        self.appendBill(bill);
                        if (allCnt === 0) {
                            continuePlease();
                        }
                    },
                    (message) => {
                        allCnt--;
                        doneCount++;
                        toast(message);
                    }
                )
            });
        }
        $("#bills-controls-btns").hide(200);
        continuePlease();
    }

    final() {
        let self = this;
        let target = this.target;
        target.year = target.period.year;
        target.month = target.period.month;
        target.closed = new Date().toISOString();
        manager.UpdateTarget(self.target_id, target, (target) => {
            self.target = target
        }, (message) => {
            toast(message);
        });
        $("#bills-controls-btns").hide(200);
    }

    destroy() {
        this.collection.unconsume(this.consumer_id);
    }
}

let draftAll = () => {};
let undraftAll = () => {};

function BillsViewPage(props) {
    let target = undefined;
    let account = undefined;
    if (props) {
        if (props.target) {
            target = props.target;
        }
        if (props.account_id) {
            account = props.account_id;
        }
    }
    let mainPage = new Render("#main-page-container");
    mainPage.content([
        {
            tag: "div", class: "row",
            content: [
                {tag: "div", id: "bills-header", class: ["col", "s12"], content: ". . ."}
            ]
        },
        {
            tag: "div", class: "row",
            content: [
                {
                    tag: "div", id: "bill-generator-buffer", class: ["col", "m5", "s12"],
                    content: [
                        {tag: "div", id: "bills-accounts-header", class: "row", content: [
                                {tag: "span", content: {tag: "b", class: ["left-align", "col", "m6"], content: "Лицевые счета"}},
                                {tag: "span", content: {tag: "a", class: ["right-align", "col", "m6"], href: "#!", content: "ВСЕ", onClick: () => {draftAll(); return false;}}}
                            ]},
                        {tag: "div", id: "bills-accounts-list", class: "collection"}
                    ]
                },
                {
                    tag: "div", id: "bill-generator-section", class: ["col", "m5", "s12"],
                    content: [
                        {tag: "div", id: "bills-prepared-header", class: "row", content: [
                                {tag: "span", content: {tag: "b", class: ["left-align", "col", "m6"], content: "Начислить"}},
                                {tag: "span", content: {tag: "a", class: ["right-align", "col", "m6"], href: "#!", content: "СБРОС", onClick: () => {undraftAll(); return false;}}}
                            ]},
                        {tag: "div", id: "bills-prepared-list", class: "collection"}
                    ]
                }
            ]
        }
    ]);
    let billsGenerator = new BillsGenerator(
        target,
        new Render("#bills-header"),
        new Render("#bills-accounts-list"),
        new Render("#bills-prepared-list"),
        accounts
    );
    draftAll = () => {
        billsGenerator.draftAll();
    }
    undraftAll = () => {
        billsGenerator.undraftAll();
    }

    manager.FindBills(account, target, (bills) => billsGenerator.fillBills(bills), (message, status) => {
        if (status === 404) {
            message = "Нет данных";
        }
        toast(message);
    })

    if (accounts.empty()) {
        accounts.loadAccounts(()=>{}, ()=>{});
    }

    return ()=>{
        billsGenerator.destroy()
    };
}
