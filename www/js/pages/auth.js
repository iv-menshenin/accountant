
function AuthPage(prop) {
    modalWindow.close();

    RenderModalWindow(
        new EditorForm("Выполните вход", {
            login: {label: "Логин", type: "text", value: "", short: true},
            password: {label: "Пароль", type: "password", value: "", short: true},
        }, (cred)=>{
            console.log(cred);
            manager.Login(cred.login, cred.password, ()=>{
                modalWindow.close();
                if (prop.back) {
                    document.location.replace("#"+prop.back)
                    return
                }
                document.location.replace("#main")
            }, (message)=>{
                toast(message);
            })
        })
    );

    return ()=>{}
}