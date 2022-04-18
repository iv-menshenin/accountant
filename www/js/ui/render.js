
class Render {
    selector = undefined;
    actors = undefined;

    renderArgs(args, animation) {
        return args.map((arg) => {
            if (isBuilder(arg)) {
                return buildHTML(arg)
            }
            return {
                content: arg,
                triggers: []
            }
        }).reduce((prev, elem) => {
            prev.content.push(elem.content);
            prev.triggers.push(...elem.triggers);
            return prev;
        }, {
            content: [],
            triggers: [],
        })
    }

    constructor(divID) {
        this.selector = $(divID);
        this.actors = $("#fixed-action-btn");
    }

    append(...args) {
        let render = this.renderArgs(args, true);
        this.selector.append(...render.content);
        if (render.triggers) {
            render.triggers.forEach((fn)=>fn())
        }
    }

    prepend(...args) {
        let render = this.renderArgs(args, true);
        this.selector.prepend(...render.content);
        if (render.triggers) {
            render.triggers.forEach((fn)=>fn())
        }
    }

    content(...args) {
        let render = this.renderArgs(args, true);
        this.selector.html(...render.content);
        if (render.triggers) {
            render.triggers.forEach((fn)=>fn())
        }
    }

    clear() {
        let self = this;
        this.selector.html("");
        this.actors.hide(200, () => {
            self.actors.html("");
        });
    }

    registerFloatingButtons(button) {
        if (!button.children) {
            button.children = [];
        }
        let tagA = {tag: "a", href: button.href, class: ["btn-floating", "btn-large", button.color], content: {tag: "i", class: ["large", "material-icons"], content: button.icon}};
        if (button.onClick) {
            tagA.onClick = button.onClick;
        }
        let floating = [];
        floating.push(tagA);
        if (button.children) {
            let children = button.children.map((btn)=>{
                return {tag: "li", content: {tag: "a", href: btn.href, class: ["btn-floating", btn.color], content: {tag: "i", class: ["material-icons"], content: btn.icon}}}
            });
            floating.push({tag: "ul", content: children});
        }
        let render = new Render("#fixed-action-btn")
        render.content(floating);
        this.actors.show(200);
    }
}

function isBuilder(obj) {
    if (!obj) {
        return false
    }
    if (Array.isArray(obj)) {
        return obj.reduce((prev, obj) => {
            return (prev === true) || isBuilder(obj);
        }, false)
    }
    if (typeof(obj) === "object") {
        let keys = Object.keys(obj);
        return keys.indexOf("tag") > -1
    }
    return false
}

function buildHTML(constructor) {
    let triggers = [];

    function buildTabFromTagOption(tagOptions) {
        if (tagOptions && tagOptions.tag) {
            let content = "";
            let postPend = "";
            if (tagOptions.content) {
                let contentObj = buildHTML(tagOptions.content);
                content = contentObj.content;
                triggers.push(...contentObj.triggers);
            }
            let keys = Object.keys(tagOptions);
            for(let nn = 0;nn < keys.length;nn++){
                let key = keys[nn];
                if (key === "tag") {
                    continue;
                }
                if (key === "afterRender") {
                    let afterRender = tagOptions[key];
                    triggers.push(afterRender);
                    continue;
                }
                if (key === "onClick") {
                    let handler = tagOptions[key];
                    let id = tagOptions.id;
                    if (!id) {
                        id = randID();
                        postPend += " id=\"" + id + "\"";
                    }
                    triggers.push(()=>$("#"+id).click(handler));
                    continue;
                }
                if (["readonly","selected","checked"].includes(key)) {
                    if (tagOptions[key]) {
                        postPend += " " + key;
                    }
                } else if (!["content"].includes(key)) {
                    let options = tagOptions[key];
                    let optionStr = "";
                    if (Array.isArray(options)) {
                        if (key === "style") {
                            optionStr = options.join(";")
                        } else {
                            optionStr = options.join(" ")
                        }
                    } else {
                        optionStr = options;
                    }
                    postPend += " " + key + "=\"" + optionStr + "\"";
                }
            }
            if (content === "" && canShortTag(tagOptions.tag)) {
                return "<" + tagOptions.tag + postPend + " />";
            }
            return "<" + tagOptions.tag + postPend + ">" + content + "</" + tagOptions.tag + ">";
        } else if (tagOptions) {
            if (typeof(tagOptions) === "string" || typeof(tagOptions) === "number") {
                return tagOptions;
            } else {
                return ""
            }
        }
    }
    function buildTags(tagOptions) {
        if (Array.isArray(tagOptions)) {
            let results = [];
            for (let nn = 0;nn < tagOptions.length;nn++) {
                results.push(buildTags(tagOptions[nn]))
            }
            return results.join("\n");
        }
        return buildTabFromTagOption(tagOptions)
    }

    if (typeof(constructor) === "object") {
        let results = [];
        if (!Array.isArray(constructor)) {
            constructor = [constructor];
        }
        for (let nn = 0;nn < constructor.length;nn++) {
            let struct = constructor[nn];
            results.push(buildTags(struct));
        }
        return {
            content: results.join("\n"),
            triggers: triggers
        }
    } else {
        if (typeof(constructor) === "string" || typeof(constructor) === "number") {
            return {
                content: constructor,
                triggers: []
            };
        } else {
            return {
                content: "",
                triggers: []
            }
        }
    }
}

function canShortTag(tag) {
    return ["hr", "br", "img", "input", "meta", "link", "source"].reduce((p, n) => p || n === tag, false);
}

function MakeCollectionPage(title, collection, builder) {
    let mainPage = new Render("#main-page-container");
    return MakeCollection(title, collection, builder, mainPage)
}

function MakeCollection(title, collection, builder, render) {
    let containerID = randID();
    let containerContent = undefined;
    let construct = {};
    if (title) {
        construct = {id: containerID, tag: "ul", class: ["collection", "with-header"], content: {tag: "li", class: "collection-header", content: {tag: "h4", content: title}}};
    } else {
        construct = {id: containerID, tag: "ul", class: ["collection"], content: []};
    }
    construct.afterRender = ()=>{
        containerContent = new Render("#"+containerID);
    };
    render.content(construct);

    let rendered = {};
    let consumer_id = collection.consume((action, element, id, href)=>{
        let listItemID = "list-item-"+id;
        if (action === "delete") {
            rendered[id] = undefined;
            $("#"+listItemID).delete();
            return
        }
        let construct = builder(element);
        let composed = [
            construct.primary,
            {tag: "span", class: "secondary-content", content: construct.secondary},
        ];
        if (rendered[id]) {
            rendered[id].content(composed);
            return;
        }
        let newItem = {tag: "a", id: listItemID, href: href, class: ["collection-item", "avatar"], content: composed};
        containerContent.append(newItem);
        rendered[id] = new Render("#"+listItemID);
    });
    return ()=>{
        collection.unconsume(consumer_id);
        containerContent.clear();
        render.clear();
    }
}

function randID() {
    return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}

function RenderModalWindow(contentManager) {
    contentManager.renderTo("#modal-frame-container", "#modal-frame-footer")
    modalWindow.open();
}