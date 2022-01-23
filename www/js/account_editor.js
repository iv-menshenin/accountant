
class accountEditorForm {
    account = {};
    LS = undefined;
    CAD = undefined;
    NAGR = undefined;
    DAGR = undefined;
    KPUR = undefined;
    DPUR = undefined;
    COMM = undefined;

    pageRender(account, options) {
        this.account = account;
        if (!options) {
            options = {};
        }
        this.LS = makeInput("Номер лицового счета", account.account, {middle: true, switch: options.switch})
        this.CAD = makeInput("Кадастровый номер", account.cad_number, {middle: true, switch: options.switch})
        this.NAGR = makeInput("Номер договора", account.agreement, {middle: true, switch: options.switch})
        this.DAGR = makeInput("Дата договора", account.agreement_date, {middle: true, switch: options.switch})
        this.KPUR = makeInput("Вид приобретения", account.purchase_kind, {middle: true, switch: options.switch})
        this.DPUR = makeInput("Дата приобретения", account.purchase_date, {middle: true, switch: options.switch}) // zeroDate = "0001-01-01T00:00:00Z"
        this.COMM = makeTextArea("Комментарий", account.comment, {switch: options.switch})
        return {
            content: [
                formHeader(getShortAddress(account)),
                {tag: "div", class: "row", content: [this.LS.content, this.CAD.content]},
                {tag: "div", class: "row", content: [this.NAGR.content, this.DAGR.content]},
                {tag: "div", class: "row", content: [this.KPUR.content, this.DPUR.content]},
                {tag: "div", class: "row", content: this.COMM.content}
            ],
            footer: [
                {tag: "div", class: "s6", content: makeButton("Сохранить изменения", {class: "action-button", onclick: accountEditor.saveAction().name})}
            ]
        };
    }

    saveAction() {
        let accountEditor = this;
        return {
            name: "accountEditor.saveAction().exec()",
            exec() {
                let value = accountEditor.LS.getValue();
                if (value) {
                    accountEditor.account.account = value;
                }
                value = accountEditor.CAD.getValue();
                if (value) {
                    accountEditor.account.cad_number = value;
                }
                value = accountEditor.NAGR.getValue();
                if (value) {
                    accountEditor.account.agreement = value;
                }
                value = accountEditor.DAGR.getValue();
                if (value) {
                    accountEditor.account.agreement_date = value;
                }
                value = accountEditor.KPUR.getValue();
                if (value) {
                    accountEditor.account.purchase_kind = value;
                }
                value = accountEditor.DPUR.getValue();
                if (value) {
                    accountEditor.account.purchase_date = value;
                }
                value = accountEditor.COMM.getValue();
                if (value) {
                    accountEditor.account.comment = value;
                }
                manager.CreateAccount(
                    accountEditor.account,
                    (account) => {
                        // todo update accounts in container
                    },
                    (error) => {
                        // todo toss error
                        console.log(error);
                        alert(error)
                    }
                )
            }
        }
    }
}

let accountEditor = new(accountEditorForm);
