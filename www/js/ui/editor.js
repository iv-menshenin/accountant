
class EditorForm {
    label = "";
    attributes = [];
    saveFn = ()=>{};
    saveButton = undefined;

    constructor(label, attributes, onSaveAction) {
        this.label = label;
        this.saveFn = onSaveAction;
        this.attributes = Object.keys(attributes).map((key) => {
            let attr = attributes[key];
            return {
                name: key,
                label: attr.label,
                type: attr.type,
                options: attr.options,
                short: (attr.short),
                value: (attr.value ? attr.value : ""),
                form: {},
            }
        })
    }

    destroy() {
        // todo
    }

    build() {
        let self = this;
        let forms = [];
        for (let i = 0; i < this.attributes.length; i++) {
            let form = this.makeForm(this.attributes[i]);
            this.attributes[i].form = form;
            forms.push(form.content);
        }
        let btnContainerID = randID();
        this.saveButton = makeButton("Сохранить", {class: "action-button save-button"}, ()=>self.saveAction());
        this.checkbox = makeFlatCheckBox(true, (e) => self.setEnabled(e));
        let b = {tag: "i", class: ["tiny", "material-icons", "red-text"], content: "block"};
        let r = {tag: "i", class: ["tiny", "material-icons", "green-text"], content: "create"};
        let switcher = {tag: "div", class: "switch", content: {tag: "label", content: [b, this.checkbox.content, {tag: "span", class: "lever"}, r]}};
        let headers = [
            {tag: "div", class: ["col", "s6"], content: formHeader(this.label)},
            {tag: "div", class: ["col", "s6", "right-align"], content: switcher},
        ];
        return {
            content: [
                {tag: "div", class: "row", content: headers},
                {tag: "div", class: "row", content: forms},
            ],
            footer: [
                {tag: "div", id: btnContainerID, class: "s6", content: this.saveButton.content}
            ],
        }
    }

    render(renderForm, renderControls) {
        let build = this.build();
        if (renderControls) {
            renderForm.content(build.content);
            renderControls.content(build.footer);
        } else {
            renderForm.content([build.content, build.footer]);
        }
    }

    renderTo(divIDForm, divIDControl) {
        let footer = undefined;
        let content = new Render(divIDForm);
        if (divIDControl) {
            footer = new Render(divIDControl);
        }
        this.render(content, footer)
    }

    disable() {
        this.setEnabled(false);
    }

    enable() {
        this.setEnabled(true);
    }

    setEnabled(enabled) {
        this.checkbox.setValue(enabled);
        this.attributes.forEach((v) => {
            if (v.form.setEnabled) {
                v.form.setEnabled(enabled);
            }
        });
        if (this.saveButton.setEnabled) {
            this.saveButton.setEnabled(enabled);
        }
    }

    makeForm(attr) {
        switch (attr.type) {
            case "text": return makeTextInput(attr.label, attr.value, {short: attr.short});
            case "number": return makeNumberInput(attr.label, attr.value, {short: attr.short});
            case "password": return makeTextInput(attr.label, attr.value, {short: attr.short, password: true});
            case "select": return makeSelect(attr.label, attr.value, {short: attr.short, options: attr.options});
            case "date": return makeDatePicker(attr.label, attr.value, {short: attr.short});
            case "checkbox": return makeCheckBox(attr.label, attr.value, {short: attr.short});
            case "multiline": return makeTextArea(attr.label, attr.value, {short: attr.short});
        }
    }

    saveAction() {
        let obj = {};
        this.attributes.forEach((attr) => {
            obj[attr.name] = attr.form.getValue();
        });
        this.saveFn(obj);
    }

}
