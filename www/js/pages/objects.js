
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
        if (this.account_id) {
            consumer(action, object, object.object_id, "#object:uuid="+object.object_id+"/account="+this.account_id);
        } else {
            consumer(action, object, object.object_id, "#object:uuid="+object.object_id);
        }
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
            toast("???? ?????????????? ??????????????????");
        });
        return
    }
    let objects = new objectsManager(options.account);
    let account = accounts.getAccount(options.account);
    if (!account) {
        account = {account: "??????"}; // todo
    }
    let destroy = MakeCollectionPage(account.account + " ??????????????", objects, buildObjectElement);
    let mainPage = new Render("#main-page-container");
    mainPage.registerFloatingButtons({href: "#object:new/account=" + options.account, icon: "add", color: "brown"})
    return ()=>{
        destroy();
        objects.onDone();
    }
}

function ObjectEditPage(options) {
    if (accounts.empty()) {
        accounts.loadAccounts(()=>{}, (message)=>{console.log(message); toast("???? ?????????????? ??????????????????")});
    }
    let editor = undefined;
    let objectInfoBlock = [
        {tag: "div", id: "object-attrs"},
        {tag: "div", id: "object-ctrls"},
    ];
    let personPage = new Render("#main-page-container");
    personPage.content(objectInfoBlock);

    let consumer_id = accounts.consume(()=>{
        let account = accounts.collection.find((account) => {
            if (options.uuid) {
                return account.objects.find((object) => object.object_id === options.uuid);
            }
            return account.account_id === options.account;
        })
        if (account) {
            let object = {}
            if (options.uuid) {
                object = account.objects.find((object) => object.object_id === options.uuid);
            }
            editor = makeObjectEditor(account.account_id, object);
            editor.renderTo("#object-attrs", "#object-ctrls");
        }
    })
    return ()=>{
        accounts.unconsume(consumer_id);
        if (editor) {
            editor.destroy();
        }
    };
}

function makeObjectEditor(account_id, object) {
    let header = getObjectShortAddress(object);
    if (!object.object_id) {
        header = "??????????";
    }
    return new EditorForm(header, {
        // postal_code: {label: "????????????", type: "text", value: object.postal_code, short: true},
        // city: {label: "??????????", type: "text", value: object.city, short: true},
        // village: {label: "???????????????????? ??????????", type: "text", value: object.village, short: true},
        street: {label: "??????????", type: "text", value: object.street, short: true},
        number: {label: "?????????? ??????????????", type: "number", value: object.number, short: true},
        area: {label: "??????????????", type: "number", value: object.area, short: false},
    }, (updated)=>{
        if (object.object_id) {
            updated.object_id = object.object_id;
            manager.UpdateObject(account_id, updated, (updatedObject)=>{
                accounts.addOrReplaceObject(account_id, updatedObject);
                toast("?????????????? ????????????????");
            }, (message)=>{
                console.log(message);
                toast("??????-???? ?????????? ???? ??????");
            })
        } else {
            manager.CreateObject(account_id, updated, (updatedObject)=>{
                accounts.addOrReplaceObject(account_id, updatedObject);
                toast("?????????????? ????????????");
                document.location.replace("#objects:account="+account_id);
            }, (message)=>{
                console.log(message);
                toast("??????-???? ?????????? ???? ??????");
            })
        }
    })
}