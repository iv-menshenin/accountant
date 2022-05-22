
function PaymentsPage(props) {
    if (props.action === "new") {
        return NewPaymentsPage(props)
    }
    if (props.uuid) {
        manager.GetPayment(props.uuid, (payment)=>{
            manager.GetAccount(payment.account_id, (account) => {
                PaymentEditorPage(account, payment)
            }, (err) => {toast(err)})
        }, (err) => {toast(err)});

        let paymentsPage = new Render("#main-page-container");
        paymentsPage.content(preloader());
        return () => {}
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

function PaymentEditorPage(account, payment) {
    let paymentInfoBlock = [
        {tag: "div", id: "payment-attrs"},
        {tag: "div", id: "payment-ctrls"},
    ];
    let paymentsPage = new Render("#main-page-container");
    paymentsPage.content(paymentInfoBlock);

    let editor = makePaymentEditor(account, payment);
    editor.renderTo("#payment-attrs", "#payment-ctrls");
}

function makePaymentEditor(account, payment) {
    let header = {tag: "span", class: ["truncate"], content: "Новая оплата " + getFirstPersonName(account)};
    if (payment) {
        header = {tag: "span", class: ["truncate"], content: "Оплата " + getFirstPersonName(account)};
    }
    if (!account.account_id) {
        return
    }
    if (!payment) {
        payment = {
            person_id: undefined,
            object_id: undefined,
            period: {
                year: (new Date()).getFullYear(),
                month: (new Date()).getMonth() + 1
            },
            target: {
                target_id: undefined
            },
            payment: 0,
            payment_date: new Date(),
            receipt: ""
        }
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
        year: {label: "Год", type: "number", value: payment.period.year, short: true},
        month: {label: "Месяц", type: "number", value: payment.period.month, short: true},
        payment: {label: "Сумма", type: "number", value: payment.payment, short: true},
        payment_date: {label: "Дата оплаты", type: "date", value: payment.payment_date, short: true},
        receipt: {label: "Номер чека", type: "text", value: payment.receipt, short: false},
        person: {label: "Владелец", type: "select", value: persons[payment.person_id], options: Object.keys(persons)},
        object: {label: "Объект", type: "select", value: objects[payment.object_id], options: Object.keys(objects)},
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