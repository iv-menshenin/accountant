
class accountsManager {
    onDone = () => {};
    collection = [];
    consumers = [];

    loadAccounts(onDone, onError) {
        this.onDone = onDone;
        manager.GetAccounts(
            (accounts) => {
                let self = this;
                this.collection = accounts;
                this.collection.forEach((account)=>{
                    self.messageAll("add", account)
                });
                onDone(accounts);
            },
            (message, status) => {
                if (status === 404) {
                    onDone([]);
                    return
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
                let personPos = persons.find((person) => person.person_id === newPerson.person_id);
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
        consumer(action, buildAccountElement(account), account.account_id, "#account:uuid="+account.account_id);
    }
}

let accounts = new(accountsManager);

function AccountsListPage() {
    if (accounts.collection.length === 0) {
        accounts.loadAccounts(()=>{}, ()=>{});
    }
    let destroy = MakeCollectionPage("Лицевые счета", accounts);
    let mainPage = new Render("#main-page-container");
    mainPage.registerFloatingButtons({href: "#account:new", icon: "add", color: "brown", onClick: () => {
        let account = {};
        RenderModalWindow(
            new EditorForm("Новый ЛС", {
                account: {label: "Номер счета", type: "text", value: account.account, short: true},
                cad_number: {label: "Кадастровый номер", type: "text", value: account.cad_number, short: true},
                agreement: {label: "Номер договора", type: "text", value: account.agreement, short: true},
                agreement_date: {label: "Дата договора", type: "date", value: account.agreement_date, short: true},
                purchase_kind: {label: "Вид собственности", type: "text", value: account.purchase_kind, short: true},
                purchase_date: {label: "Дата приобретения", type: "date", value: account.purchase_date, short: true},
                comment: {label: "Комментарий", type: "multiline", value: account.comment, short: false},
            }, (updated)=>{
                updated.account_id = account.account_id;
                manager.UpdateAccount(updated, (updatedAccount)=>{
                    accounts.addOrReplaceAccount(updatedAccount)
                }, (message)=>{
                    console.log(message);
                    toast("Что-то пошло не так");
                })
            })
        );
        return false;
    }})
    return ()=>{
        destroy();
    }
}

function AccountPage(props, retry=true) {
    if (!props.uuid) {
        return false
    }
    let account = accounts.getAccount(props.uuid);

    if (!account) {
        if (retry) {
            accounts.loadAccounts(()=> {AccountPage(props, false)}, (err) => {toast(err)});
            return false
        }
        return false
    }
    let editor = new EditorForm(getFirstPersonName(account), {
        account: {label: "Номер счета", type: "text", value: account.account, short: true},
        cad_number: {label: "Кадастровый номер", type: "text", value: account.cad_number, short: true},
        agreement: {label: "Номер договора", type: "text", value: account.agreement, short: true},
        agreement_date: {label: "Дата договора", type: "date", value: account.agreement_date, short: true},
        purchase_kind: {label: "Вид собственности", type: "text", value: account.purchase_kind, short: true},
        purchase_date: {label: "Дата приобретения", type: "date", value: account.purchase_date, short: true},
        comment: {label: "Комментарий", type: "multiline", value: account.comment, short: false},
    }, (updated)=>{
        updated.account_id = account.account_id;
        manager.UpdateAccount(updated, (updatedAccount)=>{
            accounts.addOrReplaceAccount(updatedAccount)
        }, (message)=>{
            console.log(message);
            toast("Что-то пошло не так");
        })
    });
    editor.renderTo("#main-page-container");
    return ()=>{editor.destroy()};
}