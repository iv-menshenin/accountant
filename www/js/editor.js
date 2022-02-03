
class EditorForm {
    label = "";
    attributes = [];
    render = [];
    saveFn = ()=>{};

    constructor(label, attributes, onSaveAction) {
        this.label = label;
        this.saveFn = onSaveAction;
        this.attributes = Object.keys(attributes).map((key) => {
            let attr = attributes[key];
            return {
                name: key,
                label: attr.label,
                type: attr.type,
                short: (attr.short),
                value: (attr.value ? attr.value : ""),
                form: {},
            }
        })
    }

    destroy() {

    }

    renderTo(divIDForm, divIDControl) {
        let forms = [];
        for (let i = 0; i < this.attributes.length; i++) {
            let form = this.makeForm(this.attributes[i]);
            this.attributes[i].form = form;
            forms.push(form.content);
        }
        let btnContainerID = randID();
        let formConstructor = {
            content: [
                formHeader(this.label),
                {tag: "div", class: "row", content: forms},
            ],
            footer: [
                {tag: "div", id: btnContainerID, disabled: "disabled", class: "s6", content: makeButton("Сохранить", {class: "action-button save-button"})}
            ]
        };
        let renderStruct = [formConstructor.content];
        if (divIDControl) {
            let render = new Render(divIDControl);
            render.content(formConstructor.footer);
            this.render.push(render);
        } else {
            renderStruct.push(formConstructor.footer);
        }
        let render = new Render(divIDForm);
        render.content(renderStruct);
        this.render.push(render);
        let self = this;
        $("#"+btnContainerID+" .save-button").on("click", ()=>self.saveAction());
    }

    makeForm(attr) {
        switch (attr.type) {
            case "text": return makeInput(attr.label, attr.value, {short: attr.short});
            case "password": return makeInput(attr.label, attr.value, {short: attr.short, password: true});
            case "date": return makeDatePicker(attr.label, attr.value, {short: attr.short});
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
