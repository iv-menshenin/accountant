
function getFirstPersonName(account) {
    if (account.persons && account.persons.length > 0) {
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
    return "Неизвестный";
}

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
        href: "#account:uuid=" + account.account_id,
        content: obj,
    });
}

function accountEditPageRender(account) {
    let makeInput = (id, value, label, long=false) => {
        let classes = ["input-field", "col", "m6", "s12"];
        if (long) {
            classes = ["input-field", "col", "s12"];
        }
        return {
            tag: "div", class: classes, content: [
                {id: id, tag: "input", disabled: true, type: "text", class: "validate", value: (value ? value : "")},
                {tag: "label", for: id, class: (value ? "active" : ""), content: label}
            ]
        };
    }
    return {
        tag: "div", class: "row", content: [
            {tag: "h4", content: account.account},
            {
                tag: "div", class: "row", content: [
                    makeInput("account-cad_number", account.cad_number, "Кадастровый номер", true)
                ]
            },
            {
                tag: "div", class: "row", content: [
                    makeInput("account-agreement", account.agreement, "Номер договора"),
                    makeInput("account-agreement_date", account.agreement_date, "Дата договора"),
                ]
            },
            {
                tag: "div", class: "row", content: [
                    makeInput("account-purchase_kind", account.purchase_kind, "Вид приобретения"),
                    makeInput("account-purchase_date", account.purchase_date, "Дата приобретения"),
                ]
            },
            {
                tag: "div", class: "row", content: [
                    {
                        tag: "div", class: ["input-field", "col", "s12"], content: [
                            {id: "account-comment", tag: "textarea", class: "materialize-textarea", content: account.comment},
                            {tag: "label", for: "account-comment", class: (account.comment ? "active" : ""), content: "Комментарий"}
                        ]
                    }
                ]
            },
        ]
    };
}
