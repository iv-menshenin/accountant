
// конфигурация меню навигации
navPagesList = [
    {title: "Главная", anchor: "#main", handler: ()=>{alert("main")}},
    {title: "Лицевые счета", nav: true, anchor: "#accounts", handler: (prop)=>{AccountsListPage(prop)}},
    {title: "Лицевой счет", anchor: "#account", handler: (prop)=>{AccountPage(prop)}},
    {title: "Цели", nav: true, anchor: "#targets"},
    {title: "Бухгалтерия", nav: true, anchor: "#money"},
];

function initNavigationBar() {
    for (const navPage of navPagesList) {
        // вне навигации
        if (!navPage.nav) {
            continue
        }

        // заполнение горизонтального меню навигации
        let navFullPageLine = buildHTML({
            tag: "li",
            class: "hide-on-med-and-down",
            content: tagA(navPage.title, {href: navPage.anchor}),
        });
        $("#navigation-full").append(navFullPageLine);

        // заполнение меню навигации в слайд-панели
        // для добавления горизонтального разделителя между группами меню
        // можно добавить строку типа <li><div class="divider"></div></li>
        let navSlidePageLine = buildHTML({
            tag: "li",
            content: tagA(navPage.title, {href: navPage.anchor}),
        });
        $("#nav-slide-out").append(navSlidePageLine);
    }
}

hashPattern = /(#[a-z]+):?(([a-z0-9_=-]+\/?)*)?/;

function urlHashChange(){
    let hash = location.hash;
    let hashChunks = hashPattern.exec(hash);
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
    let selectedPage = navPagesList.filter((x) => {return x.anchor === pageCode})
    if (selectedPage) {
        if (sideNav.isOpen) {
            sideNav.close()
        }
        selectedPage[0].handler(pageParameters);
        return true
    }
    console.log(pageParameters);
}
