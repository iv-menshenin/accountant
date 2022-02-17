
class accountsManager {
    onDone = () => {};
    collection = [];
    consumers = [];

    loadAccounts(onDone, onError) {
        this.onDone = onDone;
        manager.GetAccounts(
            (accounts) => {
                let self = this;
                accounts.sort((a, b)=> a.account - b.account);
                this.collection = accounts;
                this.collection.forEach((account)=>{
                    self.messageAll("add", account)
                });
                onDone(accounts);
            },
            (message, status) => {
                if (status === 404) {
                    message = "Нет данных";
                }
                toast(message);
                onError(message);
            },
        )
    }

    getAccount(account_id) {
        let filtered = this.collection.filter((account) => {return account.account_id === account_id});
        if (filtered.length === 1) {
            return filtered[0]
        }
        return undefined
    }

    addOrReplaceAccount(newAccount) {
        let self = this;
        let replaced = this.collection.reduce((replaced, account, i)=>{
            if (account.account_id === newAccount.account_id) {
                self.collection[i] = newAccount;
                return true;
            }
            return replaced;
        }, false)
        if (!replaced) {
            this.collection.push(newAccount);
        }
        this.messageAll("replace", newAccount)
        this.onDone(this.collection);
    }

    addOrReplacePerson(accountID, newPerson) {
        for (let i = 0; i < this.collection.length; i++) {
            if (this.collection[i].account_id === accountID) {
                let persons = [];
                if (this.collection[i].persons) {
                    persons = this.collection[i].persons;
                }
                let personPos = persons.findIndex((person) => person.person_id === newPerson.person_id);
                if (personPos < 0) {
                    this.collection[i].persons.push(newPerson);
                } else {
                    this.collection[i].persons[personPos] = newPerson;
                }
                this.messageAll("replace", this.collection[i]);
                return
            }
        }
    }

    consume(consumer) {
        let consumer_id = randID();
        this.consumers.push({id: consumer_id, handler: consumer});
        this.collection.forEach((account)=>{
            this.message("add", account, consumer);
        });
        return consumer_id;
    }

    unconsume(consumer_id) {
        this.consumers = this.consumers.filter((consumer)=>{return consumer.id !== consumer_id});
    }

    messageAll(action, account) {
        this.consumers.forEach((consumer) => {
            this.message(action, account, consumer.handler);
        });
    }

    message(action, account, consumer) {
        consumer(action, account, account.account_id, "#account:uuid="+account.account_id);
    }

    empty() {
        return this.collection.length === 0;
    }
}

let accounts = new(accountsManager);

function AccountsListPage() {
    let destroy = MakeCollectionPage("Лицевые счета", accounts, buildAccountElement);
    let mainPage = new Render("#main-page-container");
    mainPage.registerFloatingButtons({href: "#account:new", icon: "add", color: "brown"})
    // lazy load
    if (accounts.empty()) {
        accounts.loadAccounts(()=>{}, ()=>{});
    }
    return ()=>{
        destroy();
    }
}

function AccountPage(props, retry=true) {
    let account = {};
    let editMode = (props && props.uuid);
    if (editMode) {
        account = accounts.getAccount(props.uuid);
        if (!account) {
            if (retry) {
                accounts.loadAccounts(()=> {AccountPage(props, false)}, (err) => {toast(err)});
                return false
            }
            toast("Не удается получить данные");
            return false
        }
    }
    let editor = undefined;
    let accountInfoBlock = [
        {tag: "div", id: "account-attrs"},
        {tag: "div", id: "account-ctrls"},
    ];
    let accountPage = new Render("#main-page-container");
    if (editMode) {
        accountPage.content(makeTripleAccountBlock(account, accountInfoBlock));
    } else {
        accountPage.content(accountInfoBlock);
    }

    let consumeFn = (action, element, id) => {
        if (element.account_id !== props.uuid) {
            return
        }
        if (editor) {
            editor.destroy()
            editor = undefined;
        }
        if (action === "delete") {
            return
        }
        editor = makeAccountEditor(account);
        editor.renderTo("#account-attrs", "#account-ctrls");
        if (editMode) {
            editor.disable();
        }
    };
    let consumer_id = accounts.consume(consumeFn);
    if (!editMode) {
        consumeFn("", {}, "");
    }

    return ()=>{
        accounts.unconsume(consumer_id);
        if (editor) {
            editor.destroy();
        }
    };
}

let agreementKinds = [
    "", "договор к/п", "наследство", "постановление", "дарение"
];

function makeAccountEditor(account) {
    let header = getFirstPersonName(account);
    if (!account.account_id) {
        header = "Новый";
    }
    return new EditorForm(header, {
        account: {label: "Номер счета", type: "text", value: account.account, short: true},
        cad_number: {label: "Кадастровый номер", type: "text", value: account.cad_number, short: true},
        agreement: {label: "Номер договора", type: "text", value: account.agreement, short: true},
        agreement_date: {label: "Дата договора", type: "date", value: account.agreement_date, short: true},
        purchase_kind: {label: "Вид собственности", type: "select", options: agreementKinds, value: account.purchase_kind, short: true},
        purchase_date: {label: "Дата приобретения", type: "date", value: account.purchase_date, short: true},
        comment: {label: "Комментарий", type: "multiline", value: account.comment, short: false},
    }, (updated)=>{
        if (account.account_id) {
            updated.account_id = account.account_id;
            manager.UpdateAccount(updated, (updatedAccount)=>{
                accounts.addOrReplaceAccount(updatedAccount)
            }, (message)=>{
                console.log(message);
                toast("Что-то пошло не так");
            })
        } else {
            manager.CreateAccount(updated, (updatedAccount)=>{
                accounts.addOrReplaceAccount(updatedAccount);
                document.location.replace("#account:uuid="+updatedAccount.account_id);
            }, (message)=>{
                console.log(message);
                toast("Что-то пошло не так");
            })
        }
    })
}

function makeTripleAccountBlock(account, accountInfoBlock) {
    let collapsibleID = rndDivID();
    let personsID = rndDivID();
    let objectsID = rndDivID();
    let accHeader = (account ? account.account + "&nbsp;&nbsp;:&nbsp;:&nbsp;&nbsp;" : "") +
        getFirstPersonName(account) + "&nbsp;&nbsp;:&nbsp;:&nbsp;&nbsp;" +
        getShortAddress(account);
    let collapsibleAccount = [
        {tag: "div", class: "collapsible-header", content:
            [
                {tag: "div", class: ["truncate"], content: accHeader},
                {tag: "span", class: ["right", "badge"], content: {tag: "i", class: ["material-icons", "small", "badge"], content: "assignment"}}
            ]
        },
        {tag: "div", class: "collapsible-body", content: accountInfoBlock},
    ];
    let collapsiblePersons = [
        {tag: "div", class: ["collapsible-header", "center-block"], content:
            [
                {tag: "div", class: ["truncate"], content: "Зарегистрировано: " + personsHeader(account)},
                {tag: "span", class: ["right", "badge"], content: {tag: "a", class: ["material-icons", "small", "deep-orange-text"], href: "#persons:account="+account.account_id, content: "account_circle"}}
            ]
        },
        {tag: "div", id: personsID, class: "collapsible-body", content: []},
    ];
    let collapsibleObjects = [
        {tag: "div", class: ["collapsible-header", "center-block"], content:
            [
                {tag: "div", class: ["truncate"], content: "Участки: " + objectsHeader(account)},
                {tag: "span", class: ["right", "badge"], content: {tag: "a", class: ["material-icons", "small", "deep-orange-text"], href: "#objects:account="+account.account_id, content: "home"}}
            ]
        },
        {tag: "div", id: objectsID, class: "collapsible-body", content: []},
    ];
    return {
        tag: "ul", id: collapsibleID, class: "collapsible", content: [
            {tag: "li", content: collapsibleAccount}, // , class: "active"
            {tag: "li", content: collapsiblePersons},
            {tag: "li", content: collapsibleObjects},
        ],
        afterRender: ()=>{
            let elems = $("#" + collapsibleID);
            let instance = M.Collapsible.init(elems, {})[0];
            // free.push(() => instance.destroy());

            let pr = new Render("#"+personsID);
            let pm = new personsManager(account.account_id)
            MakeCollection("", pm, buildPersonElement, pr);

            let or = new Render("#"+objectsID);
            let om = new objectsManager(account.account_id)
            MakeCollection("", om, buildObjectElement, or);
            // todo release
        }
    };
}