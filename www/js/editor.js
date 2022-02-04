
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
        let forms = [];
        for (let i = 0; i < this.attributes.length; i++) {
            let form = this.makeForm(this.attributes[i]);
            this.attributes[i].form = form;
            forms.push(form.content);
        }
        let btnContainerID = randID();
        this.saveButton = makeButton("Сохранить", {class: "action-button save-button"});
        return {
            content: [
                formHeader(this.label),
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
        // todo following two lines you must encapsulate
        let self = this;
        $("#"+build.footer[0].id+" .save-button").on("click", ()=>self.saveAction());
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
            case "text": return makeInput(attr.label, attr.value, {short: attr.short});
            case "password": return makeInput(attr.label, attr.value, {short: attr.short, password: true});
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
