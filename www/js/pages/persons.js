
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
        consumer(action, person, person.person_id, "#person:uuid="+person.person_id+"/account="+this.account_id);
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
        accounts.loadAccounts(()=>{
            ObjectsListPage(options)
        }, (message)=>{
            console.log(message);
            toast("Не удалось загрузить");
        });
        return
    }
    let persons = new personsManager(options.account);
    let account = accounts.getAccount(options.account);
    let destroy = MakeCollectionPage(account.account + " Владельцы", persons, buildPersonElement);
    let mainPage = new Render("#main-page-container");
    mainPage.registerFloatingButtons({href: "#person:new/account:" + options.account, icon: "add", color: "brown"})
    return ()=>{
        destroy();
        persons.onDone();
    }
}