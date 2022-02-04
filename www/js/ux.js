
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

class inputForm {
    el = undefined;
    id = undefined;
    label = "";
    options = {};
    classes = [];

    constructor(options, label) {
        this.id = rndDivID();
        this.label = label;
        if (options) {
            this.options = options;
        }
        this.classes = buildClassesForInput(this.options);
    }

    content(inputType, value) {
        let self = this;
        return {
            tag: "div", class: self.classes, content: [
                {id: self.id, tag: "input", type: inputType, class: "validate", value: (value ? value : "")},
                {tag: "label", for: self.id, class: (value ? "active" : ""), content: this.label}
            ],
            afterRender: ()=> {
                let elems = $("#"+self.id);
                this.el = elems[0];
            }
        }
    }

}

function makeInput(label, value, options) {
    let form = new inputForm(options);
    let content = form.content((options.password ? "password" : "text"), value);
    return {
        content: content,
        getValue: () => {
            return form.el.value;
        },
        setEnabled: (enabled) => {
            form.el.removeAttribute("disabled");
            if (!enabled) {
                form.el.setAttribute("disabled", "true");
            }
        }
    }
}


function makeButton(content, options) {
    if (!options) {
        options = {};
    }
    let buttonID = rndDivID();
    let modal_action = undefined;
    let classes = ["btn", "waves-effect", "waves-light"];
    if (options.class) {
        classes.push(options.class);
    }
    if (options.target) {
        modal_action = options.target;
        classes.push("modal-trigger");
    }
    let instance = undefined;
    let buttonTag = {
        id: buttonID,
        tag: "button",
        class: classes,
        content: content,
        afterRender: ()=> {
            let elems = $("#"+buttonID);
            instance = elems[0];
        }
    };
    if (options.onclick) {
        buttonTag["onclick"] = options.onclick;
    }
    if (modal_action) {
        buttonTag["data-target"] = modal_action;
    }
    return {
        content: buttonTag,
        getValue: () => {},
        setEnabled: (enabled) => {
            instance.removeAttribute("disabled");
            if (!enabled) {
                instance.setAttribute("disabled", "true");
            }
        }
    }
}

function makeSelect(label, value, options) {
    if (!options) {
        options = {}
    }
    let inputID = rndDivID();
    let classes = buildClassesForInput(options);
    if (!options.options) {
        options.options = [];
    }
    let instance = undefined;
    return {
        content: {
            tag: "div", class: classes, content: [
                {tag: "label", class: "active", content: label}, // always active if there is no empty option
                {id: inputID, tag: "select", content: options.options.map(
                    (v)=>{
                        let option = {tag: "option", content: v};
                        if (v === value) {
                            option.selected = "selected";
                        }
                        return option;
                    }
                )}
            ],
            afterRender: ()=> {
                let elems = $("#"+inputID);
                instance = M.FormSelect.init(elems, options)[0];
            }
        },
        getValue: () => {
            if (instance.input) {
                return instance.input.value;
            }
        },
        setEnabled: (enabled) => {
            instance.el.removeAttribute("disabled");
            if (!enabled) {
                instance.el.setAttribute("disabled", "true");
            }
            // we need to reinitialize form
            instance.destroy();
            let elems = $("#"+inputID);
            instance = M.FormSelect.init(elems, options)[0];
        }
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
    let instance = undefined;
    return {
        content: {
            tag: "label", class: classes, content: [
                {id: textAreaID, tag: "input", type: "checkbox", checked: (value ? "checked" : "false"), content: value},
                {tag: "span", content: label}
            ],
            afterRender: ()=> {
                let elems = $("#"+textAreaID);
                instance = elems[0];
            }
        },
        getValue: () => {
            return instance.checked;
        },
        setEnabled: (enabled) => {
            instance.removeAttribute("disabled");
            if (!enabled) {
                instance.setAttribute("disabled", "true");
            }
        }
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
    let instance = undefined;
    return {
        content: {
            tag: "div", class: classes, content: [
                {id: textAreaID, tag: "textarea", class: "materialize-textarea", content: value},
                {tag: "label", for: textAreaID, class: (value ? "active" : ""), content: label}
            ],
            afterRender: ()=> {
                let elems = $("#"+textAreaID);
                instance = elems[0];
            }
        },
        getValue: () => {
            return instance.value;
        },
        setEnabled: (enabled) => {
            instance.removeAttribute("disabled");
            if (!enabled) {
                instance.setAttribute("disabled", "true");
            }
        }
    }
}

function datePickerLocalizedOptions(value) {
    return {
        defaultDate: new Date(value),
        maxDate: new Date(),
        minDate: new Date("1900-01-01T00:00:00Z"),
        setDefaultDate: true,
        yearRange: [(new Date()).getFullYear() - 100, (new Date()).getFullYear()],
        format: "dd.mm.yyyy",
        showClearBtn: true,
        showDaysInNextAndPreviousMonths: true,
        firstDay: 1,
        i18n: {
            cancel: "Отмена",
            clear: "Очистить",
            done: "Выбрать",
            months: ["Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"],
            monthsShort: ["Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"],
            weekdays: ["Воскресение", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"],
            weekdaysShort: ["Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"],
            weekdaysAbbrev: ["Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"],
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
                instance = M.Datepicker.init(elems, datePickerLocalizedOptions(value))[0];
            }
        },
        getValue: () => {
            if (instance.date) {
                // format: "yyyy-mm-ddT00:00:00Z"
                return instance.date.toISOString();
            }
        },
        setEnabled: (enabled) => {
            instance.el.removeAttribute("disabled");
            if (!enabled) {
                instance.el.setAttribute("disabled", "true");
                // we need to reinitialize form
                instance.destroy();
            }
            // todo what about re-enable after disabling?
            // let elems = $("#"+datePickerID);
            // instance = M.FormSelect.init(elems, options)[0];
        }
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