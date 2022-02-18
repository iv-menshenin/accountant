
class personsManager {
    onDone = () => {};
    collection = [];
    consumers = [];
    account_id = undefined;
    accounts_consumed = undefined;

    constructor(account_id) {
        this.account_id = account_id;
        let self = this;
        this.accounts_consumed = accounts.consume((action, element, id, href) => {
            if (account_id && id !== account_id) {
                return
            }
            if (!element.persons) {
                return
            }
            if (action === "delete") {
                element.persons.forEach((person)=> self.delete(person.person_id));
            }
            if (action === "add") {
                element.persons.forEach((person)=> self.add(person));
            }
            if (action === "replace") {
                element.persons.forEach((person)=> self.replace(person));
            }
        })
        this.onDone = ()=>{
            accounts.unconsume(self.accounts_consumed);
        }
    }

    consume(consumer) {
        let consumer_id = randID();
        this.consumers.push({id: consumer_id, handler: consumer});
        this.collection.forEach((el)=>{
            this.message("add", el, consumer);
        });
        return consumer_id;
    }

    unconsume(consumer_id) {
        this.consumers = this.consumers.filter((consumer)=>{return consumer.id !== consumer_id});
    }

    messageAll(action, person) {
        let self = this;
        this.consumers.forEach((consumer) => {
            self.message(action, person, consumer.handler);
        });
    }

    message(action, person, consumer) {
        if (this.account_id) {
            consumer(action, person, person.person_id, "#person:uuid="+person.person_id+"/account="+this.account_id);
        } else {
            consumer(action, person, person.person_id, "#person:uuid="+person.person_id);
        }
    }

    delete(person_id) {
        let self = this;
        this.collection.filter((p) => p.person_id === person_id).forEach((p) => self.messageAll("delete", p));
        this.collection = this.collection.filter((p) => p.person_id !== person_id);
    }

    add(person) {
        this.collection.push(person);
        this.messageAll("add", person)
    }

    replace(person) {
        let person_id = person.person_id;
        this.collection = this.collection.filter((p) => p.person_id !== person_id)
        this.collection.push(person);
        this.messageAll("replace", person)
    }
}

function PersonsListPage(options) {
    if (accounts.empty()) {
        accounts.loadAccounts(()=>{}, (message)=>{console.log(message); toast("Не удалось загрузить")});
    }
    let persons = new personsManager(options.account);
    let account = accounts.getAccount(options.account);
    if (!account) {
        account = {account: "Все"}; // todo
    }
    let destroy = MakeCollectionPage(account.account + " Владельцы", persons, buildPersonElement);
    let mainPage = new Render("#main-page-container");
    mainPage.registerFloatingButtons({href: "#person:new/account=" + options.account, icon: "add", color: "brown"})
    return ()=>{
        destroy();
        persons.onDone();
    }
}

function PersonEditPage(options) {
    if (accounts.empty()) {
        accounts.loadAccounts(()=>{}, (message)=>{console.log(message); toast("Не удалось загрузить")});
    }
    let editor = undefined;
    let personInfoBlock = [
        {tag: "div", id: "person-attrs"},
        {tag: "div", id: "person-ctrls"},
    ];
    let personPage = new Render("#main-page-container");
    personPage.content(personInfoBlock);

    let consumer_id = accounts.consume(()=>{
        let account = accounts.collection.find((account) => {
            if (options.uuid) {
                return account.persons.find((person) => person.person_id === options.uuid);
            }
            return account.account_id === options.account;
        })
        if (account) {
            let person = {}
            if (options.uuid) {
                person = account.persons.find((person) => person.person_id === options.uuid);
            }
            editor = makePersonEditor(account.account_id, person);
            editor.renderTo("#person-attrs", "#person-ctrls");
        }
    });
    return ()=>{
        accounts.unconsume(consumer_id);
        if (editor) {
            editor.destroy();
        }
    };
}

function makePersonEditor(account_id, person) {
    let header = getPersonFullName(person);
    if (!person.person_id) {
        header = "Новый";
    }
    return new EditorForm(header, {
        name: {label: "Имя", type: "text", value: person.name, short: true},
        surname: {label: "Фамилия", type: "text", value: person.surname, short: true},
        pat_name: {label: "Отчество", type: "text", value: person.pat_name, short: true},
        dob: {label: "Дата рождения", type: "date", value: person.dob, short: true},
        phone: {label: "Телефон", type: "text", value: person.phone, short: true},
        email: {label: "Электронная почта", type: "text", value: person.email, short: true},
        is_member: {label: "Член товарищества", type: "checkbox", value: person.is_member, short: false},
    }, (updated)=>{
        if (person.person_id) {
            updated.person_id = person.person_id;
            manager.UpdatePerson(account_id, updated, (updatedPerson)=>{
                accounts.addOrReplacePerson(account_id, updatedPerson)
            }, (message)=>{
                console.log(message);
                toast("Что-то пошло не так");
            })
        } else {
            manager.CreatePerson(account_id, updated, (updatedPerson)=>{
                accounts.addOrReplacePerson(account_id, updatedPerson)
                document.location.replace("#person:uuid="+updatedPerson.person_id);
            }, (message)=>{
                console.log(message);
                toast("Что-то пошло не так");
            })
        }
    })
}