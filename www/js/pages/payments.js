
function PaymentsPage(props) {
    if (props.action === "new") {
        return NewPaymentsPage(props)
    }
}

function NewPaymentsPage(props) {
    if (!accounts.getAccount(props.account)) {
        accounts.loadAccounts(()=> {}, (err) => {toast(err)});
    }
    let editor = undefined;
    let paymentInfoBlock = [
        {tag: "div", id: "payment-attrs"},
        {tag: "div", id: "payment-ctrls"},
    ];
    let paymentsPage = new Render("#main-page-container");
    paymentsPage.content(paymentInfoBlock);

    let consumeFn = (action, element, id) => {
        if (element.account_id !== props.account) {
            return
        }
        if (editor) {
            editor.destroy()
            editor = undefined;
        }
        if (action === "delete") {
            return
        }
        editor = makePaymentEditor(element);
        editor.renderTo("#payment-attrs", "#payment-ctrls");
    };
    let consumer_id = accounts.consume(consumeFn);

    return ()=>{
        accounts.unconsume(consumer_id);
        if (editor) {
            editor.destroy();
        }
    };
}

function makePaymentEditor(account) {
    let header = {tag: "span", class: ["truncate"], content: "Новая оплата " + getFirstPersonName(account)};
    if (!account.account_id) {
        return
    }
    let persons = account.persons.reduce((obj, person) => {
        obj[getPersonFullName(person)] = person.person_id;
        return obj;
    }, {});
    let objects = account.objects.reduce((obj, object) => {
        obj[getObjectShortAddress(object)] = object.object_id;
        return obj;
    }, {});
    return new EditorForm(header, {
        year: {label: "Год", type: "number", value: (new Date()).getFullYear(), short: true},
        month: {label: "Месяц", type: "number", value: (new Date()).getMonth() + 1, short: true},
        payment: {label: "Сумма", type: "number", value: 0, short: true},
        payment_date: {label: "Дата оплаты", type: "date", value: new Date(), short: true},
        receipt: {label: "Номер чека", type: "text", value: "", short: false},
        person: {label: "Владелец", type: "select", options: Object.keys(persons)},
        object: {label: "Объект", type: "select", options: Object.keys(objects)},
    }, (updated)=>{
        updated.person_id = persons[updated.person];
        updated.object_id = objects[updated.object];
        manager.CreatePayment(account.account_id, updated, (newPayment)=>{
            document.location.replace("#account:uuid=" + account.account_id);
            toast("Сохранено");
        }, (message)=>{
            console.log(message);
            toast("Что-то пошло не так");
        })
    })
}