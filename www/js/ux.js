

function renderAccountListElement(account) {
    let obj = [
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

function accountEditPageRender(account, options) {
    if (!options) {
        options = {};
    }
    let LS = makeInput("Номер лицового счета", account.account, {middle: true, switch: options.switch})
    let CAD = makeInput("Кадастровый номер", account.cad_number, {middle: true, switch: options.switch})
    let NAGR = makeInput("Номер договора", account.agreement, {middle: true, switch: options.switch})
    let DAGR = makeInput("Дата договора", account.agreement_date, {middle: true, switch: options.switch})
    let KPUR = makeInput("Вид приобретения", account.purchase_kind, {middle: true, switch: options.switch})
    let DPUR = makeInput("Дата приобретения", account.purchase_date, {middle: true, switch: options.switch}) // zeroDate = "0001-01-01T00:00:00Z"
    let COMM = makeTextArea("Комментарий", account.comment, {switch: options.switch})
    return {
        tag: "div", class: "row", content: [
            {tag: "h4", content: account.account},
            {tag: "div", class: "row", content: [LS.content, CAD.content]},
            {tag: "div", class: "row", content: [NAGR.content, DAGR.content]},
            {tag: "div", class: "row", content: [KPUR.content, DPUR.content]},
            {tag: "div", class: "row", content: COMM.content},
        ]
    };
}

function makeButton(content, options) {
    let modal_action = "modal-action";
    let classes = ["btn", "waves-effect", "waves-light", "modal-trigger"];
    if (options && options.class) {
        classes.push(options.class);
    }
    if (options && options.target) {
        modal_action = options.target;
    }
    return {
        tag: "button",
        "data-target": modal_action,
        class: classes,
        content: content
    }
}

function makeCollectionContainer(header, options) {
    let containerID = rndDivID();
    let classes = ["collection"];
    let content = "";
    if (header !== "") {
        classes.push("with-header");
        content = {tag: "li", class: "collection-header", content: {tag: "h4", content: header}}
    }
    if (options && options.class) {
        classes.push(options.class);
    }
    return {
        content: {tag: "div", class: "", content: {id: containerID, tag: "ul", class: classes, content: content}},
        append: (content)=>{
            $("#"+containerID).append(buildHTML(content));
        },
        clear: ()=>{
            $("#"+containerID).html("");
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
            return $("#"+inputID).value;
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
            return $("#"+textAreaID).value;
        },
    }
}

function rndDivID() {
    return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}