
function buildTargetElement(target) {
    return {
        primary: [
            {tag: "img", class: "circle", src: "png/badge_catalog.png"},
            {tag: "span", class: ["title", "black-text"], content: target.type},
            {tag: "p", class: ["grey-text"], content: target.comment},
        ],
        secondary: (target.period.month < 10 ? "0" : "") + target.period.month + "." + target.period.year
    };
}

function targetHeader(target) {
    if (target.period) {
        return (target.period.month < 10 ? "0" : "") + target.period.month + "." + target.period.year + ": " + target.type
    }
    return "Новый взнос"
}

function buildAccountElement(account) {
    return {
        primary: [
            {tag: "img", class: "circle", src: "png/badge_account.png"},
            {tag: "span", class: ["title", "black-text", "truncate"], content: "<b>" + account.account + "</b><br />" + getAllPersonNames(account)}, // ФИО
            {tag: "p", class: ["grey-text"], content: getShortAddress(account)},     // информация об участках
        ],
        secondary: ""
    };
}

function buildPersonElement(person) {
    return {
        primary: [
            {tag: "img", class: "circle", src: "png/butterfly.png"},
            {tag: "span", class: ["title", "black-text", "truncate"], content: getPersonFullName(person)}, // ФИО
            {tag: "p", class: ["grey-text"], content: (person.phone ? person.phone : "")},     // телефон
        ],
        secondary: (person.is_member ? "член" : "не член")
    };
}

function buildObjectElement(object) {
    return {
        primary: [
            {tag: "img", class: "circle", src: "png/badge_object.png"},
            {tag: "span", class: ["title", "black-text", "truncate"], content: getObjectShortAddress(object)},
            {tag: "p", class: ["grey-text"], content: (object.village ? object.village : "") + " " + (object.area ? "[" + object.area + " m2]" : "?")},
        ],
        secondary: ""
    };
}

function buildBillElement(bill) {
    return {
        primary: [
            {tag: "img", class: "circle", src: "png/bills.png"},
            {tag: "span", class: ["title", "black-text", "truncate"], content: bill.bill + " руб."},
            {tag: "p", class: ["grey-text"], content: bill.target.type + " " + (new Date(bill.formed).toLocaleDateString())},
        ],
        secondary: ""
    };
}

function buildPaymentElement(payment) {
    return {
        primary: [
            {tag: "img", class: "circle", src: "png/money.png"},
            {tag: "span", class: ["title", "black-text", "truncate"], content: payment.payment + " руб. (" + payment.target.type + ")"},
            {tag: "p", class: ["grey-text", "truncate"], content: (new Date(payment.payment_date).toLocaleDateString() + " " + payment.receipt)},
        ],
        secondary: ""
    };
}

const noOwner = "Нет владельца";
const noObjects = "Без участка";

function getFirstPersonName(account) {
    if (account && account.persons && account.persons.length > 0) {
        let result = [];
        if (account.persons[0].surname) {
            result.push(account.persons[0].surname);
        }
        if (account.persons[0].name) {
            result.push(account.persons[0].name);
        }
        if (account.persons[0].pat_name) {
            result.push(account.persons[0].pat_name);
        }
        return result.join(" ");
    }
    return noOwner;
}

function getAllPersonNames(account) {
    if (isPersonPresent(account)) {
        return account.persons.map(getPersonFullName).join("; ");
    }
    return noOwner;
}

function personsHeader(account) {
    if (isPersonPresent(account)) {
        let owners = account.persons.map(getPersonFullName).join("; ");
        if (account.persons.length > 1) {
            owners = "(" + account.persons.length + ") " + owners;
        }
        return owners;
    }
    return noOwner;
}

function accountHeader(account) {
    if (!account) {
        return "None";
    }
    return account.account;
}

function getPersonFullName(person) {
    if (!person) {
        return ""
    }
    let result = [];
    if (person.surname) {
        result.push(person.surname);
    }
    if (person.name) {
        result.push(person.name);
    }
    if (person.pat_name) {
        result.push(person.pat_name);
    }
    return result.join(" ");
}

function getShortAddress(account) {
    if (isObjectPresent(account)) {
        return concatObjectsAddr(account, getObjectShortAddress);
    }
    return noObjects;
}

function objectsHeader(account) {
    if (isObjectPresent(account)) {
        let objects = concatObjectsAddr(account, getObjectShortAddress);;
        if (account.objects.length > 1) {
            objects = "(" + account.objects.length + ") " + objects;
        }
        return objects
    }
    return noObjects;
}

function isObjectPresent(account) {
    return (account && account.objects && account.objects.length > 0);
}

function isPersonPresent(account) {
    return (account && account.persons && account.persons.length > 0);
}

function concatObjectsAddr(account, fn) {
    return account.objects.map(fn).join("; ");
}

function getObjectShortAddress(object) {
    if (!object) {
        return ""
    }
    let result = [];
    // if (object.village) {
    //     result.push(object.village);
    // }
    if (object.street) {
        result.push(object.street);
    }
    if (object.number) {
        result.push(object.number);
    }
    return result.join(", ");
}
