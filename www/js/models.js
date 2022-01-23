
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
        return account.persons.map(getPersonFullName).join("; ");
    }
    return "Не зарегистрировано";
}

function getPersonFullName(person) {
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
    if (account.objects && account.objects.length > 0) {
        account.objects.map(getObjectShortAddress).join("; ");
    }
    return "Без участка";
}

function getObjectShortAddress(object) {
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
}
