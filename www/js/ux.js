

function mapAccountToListElement(account) {
    let obj = [
        {tag: "img", class: "circle", src: "/www/png/badge_account.png"},
        {tag: "span", class: ["title", "black-text"], content: getAllPersonNames(account)}, // ФИО
        {tag: "p", class: ["grey-text"], content: getShortAddress(account)},     // информация об участках
        {tag: "p", class: ["secondary-content"], content: account.account},      // лицевой счет
    ];
    return {
        tag: "a",
        class: ["collection-item", "avatar"],
        href: "#account:uuid=" + account.account_id,
        content: obj,
    };
}

function mapPersonToListElement(person) {
    let obj = [
        {tag: "img", class: "circle", src: "/www/png/butterfly.png"},
        {tag: "span", class: ["title", "black-text"], content: getPersonFullName(person)}, // ФИО
        {tag: "p", class: ["grey-text"], content: (person.phone ? person.phone : "")},     // телефон
        {tag: "p", class: ["secondary-content"], content: (person.is_member ? "член" : "не член")},
    ];
    return {
        tag: "a",
        class: ["collection-item", "avatar"],
        href: "#person:uuid=" + person.account_id,
        content: obj,
    };
}

function mapObjectToListElement(object) {
    let obj = [
        {tag: "img", class: "circle", src: "/www/png/badge_object.png"},
        {tag: "span", class: ["title", "black-text"], content: getObjectShortAddress(object)},
        {tag: "p", class: ["grey-text"], content: (object.city ? object.city : "")},
        {tag: "p", class: ["secondary-content"], content: (object.area ? object.area : "?")},
    ];
    return {
        tag: "a",
        class: ["collection-item", "avatar"],
        href: "#object:uuid=" + object.account_id,
        content: obj,
    };
}

function makeButton(content, options) {
    if (!options) {
        options = {};
    }
    let modal_action = undefined;
    let classes = ["btn", "waves-effect", "waves-light"];
    if (options.class) {
        classes.push(options.class);
    }
    if (options.target) {
        modal_action = options.target;
        classes.push("modal-trigger");
    }
    let buttonTag = {
        tag: "button",
        class: classes,
        content: content
    };
    if (options.onclick) {
        buttonTag["onclick"] = options.onclick;
    }
    if (modal_action) {
        buttonTag["data-target"] = modal_action;
    }
    return buttonTag
}

function makeCollectionContainer(header, options) {
    let containerID = rndDivID();
    let classes = ["collection"];
    let defaultContent = "";
    if (header !== "") {
        classes.push("with-header");
        defaultContent = {tag: "li", class: "collection-header", content: {tag: "h4", content: header}}
    }
    if (options && options.class) {
        classes.push(options.class);
    }
    return {
        content: {tag: "div", class: "", content: {id: containerID, tag: "ul", class: classes, content: defaultContent}},
        append: (content)=>{
            $("#"+containerID).append(buildHTML(content));
        },
        clear: ()=>{
            $("#"+containerID).html(buildHTML(defaultContent));
        },
    }
}

function makeInput(label, value, options) {
    let inputID = rndDivID();
    let classes = ["input-field", "col", "s12"];
    if (options && options.short) {
        classes.push("m6")
    }
    if (options && options.middle) {
        classes.push("l6")
    }
    if (options && options.switch) {
        options.switch((enabled)=>{
            let div = $("#"+inputID);
            div.removeAttr("disabled");
            if (!enabled) {
                div.attr("disabled", "true");
            }
        })
    }
    return {
        content: {
            tag: "div", class: classes, content: [
                {id: inputID, tag: "input", disabled: true, type: "text", class: "validate", value: (value ? value : "")},
                {tag: "label", for: inputID, class: (value ? "active" : ""), content: label}
            ]
        },
        getValue: () => {
            return $("#"+inputID)[0].value;
        },
    }
}

function makeTextArea(label, value, options) {
    let textAreaID = rndDivID();
    let classes = ["input-field", "col", "s12"];
    if (options && options.short) {
        classes.push("m6")
    }
    if (options && options.middle) {
        classes.push("l6")
    }
    if (options && options.switch) {
        options.switch((enabled)=>{
            let div = $("#"+textAreaID);
            div.removeAttr("disabled");
            if (!enabled) {
                div.attr("disabled", "true");
            }
        })
    }
    return {
        content: {
            tag: "div", class: classes, content: [
                {id: textAreaID, tag: "textarea", disabled: true, class: "materialize-textarea", content: value},
                {tag: "label", for: textAreaID, class: (value ? "active" : ""), content: label}
            ]
        },
        getValue: () => {
            return $("#"+textAreaID)[0].value;
        },
    }
}

function rndDivID() {
    return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}