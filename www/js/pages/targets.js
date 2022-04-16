
class targetsManager {
    onDone = () => {};
    collection = [];
    consumers = [];

    loadTargets(onDone, onError) {
        this.onDone = onDone;
        manager.FindTargets(
            true,
            undefined,
            undefined,
            (targets) => {
                let self = this;
                // descending sort
                targets.sort((a, b)=> (b.period.month + b.period.year * 12) - (a.period.month + a.period.year * 12));
                this.collection = targets;
                this.collection.forEach((target)=>{
                    self.messageAll("add", target)
                });
                onDone(targets);
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

    getTarget(target_id) {
        let filtered = this.collection.filter((target) => {return target.target_id === target_id});
        if (filtered.length === 1) {
            return filtered[0]
        }
        return undefined
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

    message(action, target, consumer) {
        consumer(action, target, target.target_id, "#target:uuid="+target.target_id);
    }

    empty() {
        return this.collection.length === 0;
    }
}

let targets = new(targetsManager);

function TargetsListPage() {
    let destroy = MakeCollectionPage("Взносы (целевые и ежегодные)", targets, buildTargetElement);
    let mainPage = new Render("#main-page-container");
    mainPage.registerFloatingButtons({href: "#target:new", icon: "add", color: "brown"})
    // lazy load
    if (targets.empty()) {
        targets.loadTargets(()=>{}, ()=>{});
    }
    return ()=>{
        destroy();
    }
}

function TargetViewPage(props, retry=true) {
    let target = {
        period: {
            year: (new Date()).getFullYear(),
            month: (new Date()).getMonth() + 1
        }
    };
    let editMode = (props && props.uuid);
    if (editMode) {
        target = targets.getTarget(props.uuid);
        if (!target) {
            if (retry) {
                targets.loadTargets(()=> {TargetViewPage(props, false)}, (err) => {toast(err)});
                return false
            }
            toast("Не удается получить данные");
            return false
        }
    }
    let editor = undefined;
    let accountInfoBlock = [
        {tag: "div", id: "target-attrs"},
        {tag: "div", id: "target-ctrls"},
    ];
    let targetPage = new Render("#main-page-container");
    targetPage.content(accountInfoBlock);

    let consumeFn = (action, element, id) => {
        if (element.target_id !== props.uuid) {
            return
        }
        if (editor) {
            editor.destroy()
            editor = undefined;
        }
        if (action === "delete") {
            return
        }
        editor = makeTargetEditor(target);
        editor.renderTo("#target-attrs", "#target-ctrls");
        if (editMode) {
            editor.disable();
        }
    };
    let consumer_id = targets.consume(consumeFn);
    if (!editMode) {
        consumeFn("", {}, "");
    }

    return ()=>{
        targets.unconsume(consumer_id);
        if (editor) {
            editor.destroy();
        }
    };
}

function makeTargetEditor(target) {
    let header = targetHeader(target);
    if (!target.target_id) {
        header = "Новый";
    }
    return new EditorForm(header, {
        year: {label: "Год", type: "number", value: target.period.year, short: true},
        month: {label: "Месяц", type: "number", value: target.period.month, short: true},
        type: {label: "Тип сбора", type: "text", value: target.type, short: true},
        cost: {label: "Сумма сбора", type: "number", value: target.cost, short: true},
        comment: {label: "Описание", type: "text", value: target.comment, short: false},
    }, (updated)=>{
        if (target.target_id) {
            // updated.target_id = target.target_id;
            // manager.UpdateAccount(updated, (updatedAccount)=>{
            //     accounts.addOrReplaceAccount(updatedAccount);
            //     toast("Владелец обновлен");
            // }, (message)=>{
            //     console.log(message);
            //     toast("Что-то пошло не так");
            // })
        } else {
            manager.CreateTarget(updated, (updatedTarget)=>{
                targets.addOrReplaceTarget(updatedTarget);
                toast("Новый сбор добавлен");
                document.location.replace("#target:uuid="+updatedTarget.target_id);
            }, (message)=>{
                console.log(message);
                toast("Что-то пошло не так");
            })
        }
    })
}