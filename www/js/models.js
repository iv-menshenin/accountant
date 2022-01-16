
function getAllPersonNames(account) {
    if (account.persons && account.persons.length > 0) {
        account.persons.map((person) => {
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
        }).join("; ");
    }
    return "Не зарегистрировано";
}

function getShortAddress(account) {
    if (account.objects && account.objects.length > 0) {
        account.objects.map((object) => {
            let result = [];
            if (object.village) {
                result.push(object.village);
            }
            if (object.street) {
                result.push(object.street);
            }
            if (object.number) {
                result.push(object.number);
            }
            return result.join(", ");
        }).join("; ");
    }
    return "Без участка";
}

function renderAccountListElement(account) {
    let obj = [
        {tag: "span", class: ["title", "black-text"], content: getAllPersonNames(account)}, // ФИО
        {tag: "p", class: ["grey-text"], content: getShortAddress(account)},     // информация об участках
        {tag: "p", class: ["secondary-content"], content: account.account},      // лицевой счет
    ];
    return buildHTML({
        tag: "a",
        class: ["collection-item", "avatar"],
        href: "#account:" + account.account_id,
        onclick: "showAccountDetails(this)",
        content: obj,
    });
}
