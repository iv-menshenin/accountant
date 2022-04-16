
class navigationManager {
    sideNav = undefined;
    navigationBar = undefined;
    navigationSlider = undefined;
    hashPattern = /(#[a-z]+):?(([a-z0-9_=-]+\/?)*)?/;
    destructorFn = ()=>{};
    navPagesList = [
        // главная
        {title: "Главная", anchor: "#main", handler: ()=>{return ()=>{};}},
        // реагируют на якорь
        {title: "Вход", anchor: "#login", handler: (prop)=>{return AuthPage(prop)}},
        {title: "Лицевой счет", anchor: "#account", handler: (prop)=>{return AccountPage(prop)}},
        {title: "Проживающий", anchor: "#person", handler: (prop)=>{return PersonEditPage(prop)}},
        {title: "Участок", anchor: "#object", handler: (prop)=>{return ObjectEditPage(prop)}},
        {title: "Взносы", anchor: "#target", handler: (prop)=>{return TargetViewPage(prop)}},
        // компоненты меню
        {title: "Лицевые счета", nav: true, anchor: "#accounts", handler: (prop)=>{return AccountsListPage(prop)}},
        {title: "Владельцы", nav: true, anchor: "#persons", handler: (prop)=>{return PersonsListPage(prop)}},
        {title: "Участки", nav: true, anchor: "#objects", handler: (prop)=>{return ObjectsListPage(prop)}},
        {title: "Взносы", nav: true, anchor: "#targets", handler: (prop)=>{return TargetsListPage(prop)}},
    ];

    constructor(sideNav) {
        this.sideNav = sideNav;
        this.navigationBar = new Render("#navigation-full");
        this.navigationSlider = new Render("#nav-slide-out");
        for (const navPage of this.navPagesList) {
            // вне навигации
            if (!navPage.nav) {
                continue
            }

            // заполнение горизонтального меню навигации
            this.navigationBar.append({tag: "li", class: "hide-on-med-and-down", content: tagA(navPage.title, {href: navPage.anchor})});

            // заполнение меню навигации в слайд-панели
            // для добавления горизонтального разделителя между группами меню
            // можно добавить строку типа <li><div class="divider"></div></li>
            this.navigationSlider.append({tag: "li", content: tagA(navPage.title, {href: navPage.anchor})});
        }
    }

    urlHashChange() {
        let hashChunks = this.hashPattern.exec(location.hash);
        if (!hashChunks || !hashChunks[1]) {
            hashChunks = ["#main", "#main"];
        }
        let pageCode = hashChunks[1];
        let pageParameters = {};
        if (hashChunks[2]) {
            hashChunks[2].split("/").forEach((parameter)=>{
                let p = parameter.split("=");
                if (p.length > 1) {
                    pageParameters[p[0]] = p[1];
                } else {
                    pageParameters[p[0]] = true;
                }
            });
        }
        if (this.destructorFn) {
            this.destructorFn();
        }
        let selectedPage = this.navPagesList.filter((x) => {return x.anchor === pageCode})
        if (selectedPage.length > 0) {
            if (this.sideNav.isOpen) {
                this.sideNav.close()
            }
            this.destructorFn = selectedPage[0].handler(pageParameters);
            return true
        }
        console.log(pageParameters);
    }
}
