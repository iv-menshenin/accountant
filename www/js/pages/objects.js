
class objectsManager {
    onDone = () => {};
    collection = [];
    consumers = [];
    account_id = undefined;

    constructor(account_id) {
        this.account_id = account_id;
        let self = this;
        this.accounts_consumed = accounts.consume((action, element, id, href) => {
            if (account_id && id !== account_id) {
                return
            }
            if (!element.objects) {
                return
            }
            if (action === "delete") {
                element.objects.forEach((object)=> self.delete(object.object_id));
            }
            if (action === "add") {
                element.objects.forEach((object)=> self.add(object));
            }
            if (action === "replace") {
                element.objects.forEach((object)=> self.replace(object));
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

    messageAll(action, object) {
        let self = this;
        this.consumers.forEach((consumer) => {
            self.message(action, object, consumer.handler);
        });
    }

    message(action, object, consumer) {
        consumer(action, object, object.object_id, "#object:uuid="+object.object_id+"/account="+this.account_id);
    }

    delete(object_id) {
        let self = this;
        this.collection.filter((o) => o.object_id === object_id).forEach((o) => self.messageAll("delete", o));
        this.collection = this.collection.filter((o) => o.object_id !== object_id);
    }

    add(object) {
        this.collection.push(object);
        this.messageAll("add", object)
    }

    replace(object) {
        let object_id = object.object_id;
        this.collection = this.collection.filter((o) => o.object_id !== object_id)
        this.collection.push(object);
        this.messageAll("replace", object)
    }
}

function ObjectsListPage(options) {
    if (accounts.empty()) {
        accounts.loadAccounts(()=>{
            ObjectsListPage(options)
        }, (message)=>{
            console.log(message);
            toast("Не удалось загрузить");
        });
        return
    }
    let objects = new objectsManager(options.account);
    let account = accounts.getAccount(options.account);
    let destroy = MakeCollectionPage(account.account + " Участки", objects, buildObjectElement);
    let mainPage = new Render("#main-page-container");
    mainPage.registerFloatingButtons({href: "#object:new/account:" + options.account, icon: "add", color: "brown"})
    return ()=>{
        destroy();
        objects.onDone();
    }
}