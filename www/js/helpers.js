const animDurationDefault = 200;

function preparePage(title, constructor, onComplete) {
    let trigger = ()=>{};
    let hidden = false;
    $("#main-page-container").hide({
        duration: animDurationDefault,
        complete: ()=>{
            hidden = true;
            $("#page-title").html(title);
            $("#main-page-container").html(buildHTML(constructor));
            trigger();
        },
    });
    return ()=>{
        trigger = ()=>{
            $("#main-page-container").show(animDurationDefault, onComplete);
        }
        if (hidden) {
            trigger();
        }
    }
}

function prepareModalForm(constructor) {
    $("#modal-frame-container").hide().html(buildHTML(constructor)).show();
}

function buildHTML(constructor){
    function buildTag(tagOptions){
        if (tagOptions && tagOptions.tag) {
            let content = "";
            let tag = "";
            let postPend = "";
            if (tagOptions.content) {
                content = buildHTML(tagOptions.content);
            }
            let keys = Object.keys(tagOptions);
            for(let nn = 0;nn < keys.length;nn++){
                let key = keys[nn];
                if (key === "tag") {
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
            if (content === "") {
                tag = "<" + tagOptions.tag + postPend + " />";
            } else {
                tag = "<" + tagOptions.tag + postPend + ">" + content + "</" + tagOptions.tag + ">";
            }
            return tag;
        } else if (tagOptions) {
            if (typeof(tagOptions) === "string" || typeof(tagOptions) === "number") {
                return tagOptions;
            } else {
                return ""
            }
        }
    }
    if (typeof(constructor) === "object") {
        let results = [];
        if (!Array.isArray(constructor)) {
            constructor = [constructor];
        }
        for (let nn = 0;nn < constructor.length;nn++) {
            results.push(buildTag(constructor[nn]))
        }
        return results.join("\n");
    } else {
        if (typeof(constructor) === "string" || typeof(constructor) === "number") {
            return constructor;
        } else {
            return ""
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