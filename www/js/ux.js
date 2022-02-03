

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

function buildClassesForInput(options) {
    let classes = ["input-field", "col", "s12"];
    if (options.short) {
        classes.push("l6");
    }
    if (options.class) {
        classes.push(options.class);
    }
    return classes;
}

function makeInput(label, value, options) {
    if (!options) {
        options = {}
    }
    let inputType = "text";
    if (options.password) {
        inputType = "password"
    }
    let inputID = rndDivID();
    let classes = buildClassesForInput(options);
    if (options.switch) {
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
                {id: inputID, tag: "input", type: inputType, class: "validate", value: (value ? value : "")},
                {tag: "label", for: inputID, class: (value ? "active" : ""), content: label}
            ]
        },
        getValue: () => {
            return $("#"+inputID)[0].value;
        },
    }
}

function makeCheckBox(label, value, options) {
    if (!options) {
        options = {}
    }
    let textAreaID = rndDivID();
    let classes = buildClassesForInput(options);
    if (options.switch) {
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
            tag: "label", class: classes, content: [
                {id: textAreaID, tag: "input", type: "checkbox", checked: (value ? "checked" : "false"), content: value},
                {tag: "span", content: label}
            ]
        },
        getValue: () => {
            return $("#"+textAreaID)[0].checked;
        },
    }
}

function makeTextArea(label, value, options) {
    if (!options) {
        options = {}
    }
    let textAreaID = rndDivID();
    let classes = buildClassesForInput(options);
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
                {id: textAreaID, tag: "textarea", class: "materialize-textarea", content: value},
                {tag: "label", for: textAreaID, class: (value ? "active" : ""), content: label}
            ]
        },
        getValue: () => {
            return $("#"+textAreaID)[0].value;
        },
    }
}

function makeDatePicker(label, value, options) {
    if (!options) {
        options = {}
    }
    let instance = undefined;
    let datePickerID = rndDivID();
    let classes = buildClassesForInput(options);
    if (options && options.switch) {
        options.switch((enabled)=>{
            let div = $("#"+datePickerID);
            div.removeAttr("disabled");
            if (!enabled) {
                div.attr("disabled", "true");
            }
        })
    }
    return {
        content: {
            tag: "div", class: classes, content: [
                {id: datePickerID, tag: "input", type: "text", class: "datepicker", content: ""},
                {tag: "label", for: datePickerID, class: (value ? "active" : ""), content: label}
            ],
            afterRender: ()=>{
                let elems = $("#"+datePickerID);
                instance = M.Datepicker.init(elems, {
                    defaultDate: new Date(value),
                    setDefaultDate: true,
                })[0];
            }
        },
        getValue: () => {
            // format: "yyyy-mm-ddT00:00:00Z"
            return instance.date.toISOString();
        },
    }
}

function tagA(content, options) {
    if (!options) {
        options = {};
    }
    options.tag = "a";
    options.content = content;
    return options;
}

function formHeader(title) {
    return {tag: "h5", content: title};
}

function windowHeader(title) {
    return {tag: "h4", content: title};
}