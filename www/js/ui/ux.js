
class inputForm {
    el = undefined;
    id = undefined;
    label = "";
    options = {};
    classes = [];
    type = "text";

    constructor(label, options) {
        this.id = rndDivID();
        this.label = label;
        if (options) {
            this.options = options;
        }
        this.classes = buildClassesForInput(this.options);
    }

    init(value) {
        let elems = $("#" + this.id);
        switch (this.type) {
            case "select":
                // todo add it to destroy?
                this.el = M.FormSelect.init(elems, this.options)[0];
                break;
            case "datepicker":
                // todo add it to destroy?
                this.el = M.Datepicker.init(elems, datePickerLocalizedOptions(value))[0];
                break;
            default:
                this.el = elems[0];
        }
    }

    formTag(value) {
        // todo simplify
        let area = (this.type === "textarea");
        let select = (this.type === "select");
        let dt = (this.type === "datepicker");
        let tag = (area ? "textarea" : (select ? "select" : "input"));
        let cl = (area ? "materialize-textarea" : (dt ? "datepicker" : "validate"));
        let tp = (dt ? "text" : this.type);
        let construct = {id: this.id, tag: tag, type: tp, class: cl};
        switch (this.type) {
            case "select":
                construct.content = this.options.options.map((v)=>({tag: "option", content: v, selected: (v === value ? "selected" : undefined)}));
                break;
            case "textarea":
                construct.content = (value ? value : "");
                break;
            case "datepicker":
                // todo set value here
                break;
            case "checkbox":
                construct.checked = (value ? "checked" : undefined);
                break;
            default:
                construct.value = (value ? value : "");
        }
        return construct;
    }

    labelTag(value, formTag) {
        if (this.type === "checkbox") {
            return {tag: "label", content: [formTag, {tag: "span", content: this.label}]}
        }
        return {
            tag: "div", class: this.classes, content: [
                {tag: "label", for: this.id, class: (value || this.type === "select" ? "active" : ""), content: this.label},
                formTag,
            ],
        }
    }

    content(inputType, value) {
        this.type = inputType;
        let self = this;
        let construct = this.labelTag(value, this.formTag(value))
        construct.afterRender = ()=>self.init(value);
        return construct;
    }
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

function makeTextInput(label, value, options) {
    let form = new inputForm(label, options);
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

function makeNumberInput(label, value, options) {
    let form = new inputForm(label, options);
    let content = form.content("number", value);
    return {
        content: content,
        getValue: () => {
            return parseInt(form.el.value);
        },
        setEnabled: (enabled) => {
            form.el.removeAttribute("disabled");
            if (!enabled) {
                form.el.setAttribute("disabled", "true");
            }
        }
    }
}

function makeTextArea(label, value, options) {
    let form = new inputForm(label, options);
    let content = form.content("textarea", value);
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

function makeSelect(label, value, options) {
    let form = new inputForm(label, options);
    let content = form.content("select", value);
    return {
        content: content,
        getValue: () => {
            if (form.el.input) {
                return form.el.input.value;
            }
        },
        setEnabled: (enabled) => {
            form.el.el.removeAttribute("disabled");
            if (!enabled) {
                form.el.el.setAttribute("disabled", "true");
            }
            // we need to reinitialize form
            form.el.destroy();
            form.init();
        }
    }
}

function makeDatePicker(label, value, options) {
    let form = new inputForm(label, options);
    let content = form.content("datepicker", value);
    return {
        content: content,
        getValue: () => {
            if (form.el.date) {
                // format: "yyyy-mm-ddT00:00:00Z"
                return form.el.date.toISOString();
            }
        },
        setEnabled: (enabled) => {
            form.el.el.removeAttribute("disabled");
            if (!enabled) {
                form.el.el.setAttribute("disabled", "true");
                // we need to reinitialize form
                form.el.destroy();
            }
            form.init(value);
        }
    }
}

function makeButton(content, options, onClick) {
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
            if (onClick) {
                elems.on("click", onClick);
            }
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

function makeCheckBox(label, value, options) {
    let form = new inputForm(label, options);
    let content = form.content("checkbox", value);
    return {
        content: content,
        getValue: () => {
            return form.el.checked;
        },
        setValue: (value) => {
            form.el.checked = value;
        },
        setEnabled: (enabled) => {
            form.el.removeAttribute("disabled");
            if (!enabled) {
                form.el.setAttribute("disabled", "true");
            }
        }
    }
}

function makeFlatCheckBox(value, onChange) {
    let id = rndDivID();
    let instance = undefined;
    return {
        content: {
            id: id, tag: "input", type: "checkbox", checked: (value ? "checked" : undefined),
            afterRender: ()=> {
                let elems = $("#"+id);
                instance = elems[0];
                elems.on("change", ()=>onChange(instance.checked));
            }
        },
        getValue: () => {
            return instance.checked;
        },
        setValue: (value) => {
            instance.checked = value;
        },
        setEnabled: (enabled) => {
            instance.removeAttribute("disabled");
            if (!enabled) {
                instance.setAttribute("disabled", "true");
            }
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

function preloader() {
    let preloader = {
        tag: "div",
        class: ["preloader-wrapper", "big", "active"],
        content: {
            tag: "div",
            class: ["spinner-layer", "spinner-blue-only"],
            content: [
                {
                    tag: "div",
                    class: ["circle-clipper", "left"],
                    content: {tag: "div", class: "circle"}
                },
                {
                    tag: "div",
                    class: ["gap-patch"],
                    content: {tag: "div", class: "circle"}
                },
                {
                    tag: "div",
                    class: ["circle-clipper", "right"],
                    content: {tag: "div", class: "circle"}
                },
            ]
        }
    };
    return {
        tag: "div", class: ["row", "valign-wrapper"], style: "height: 60%",
        content: {
            tag: "div", class: ["col", "s12", "center"],
            content: preloader
        }
    }
}