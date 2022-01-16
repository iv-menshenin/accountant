
// конфигурация меню навигации
function initNavigationBar() {
    let navPages = [
        {title: "Лицевые счета", anchor: "accounts", onClick: ()=>{alert("1");}},
        {title: "Цели", anchor: "targets", onClick: ()=>{alert("2");}},
        {title: "Бухгалтерия", anchor: "money", onClick: ()=>{alert("3");}},
    ];
    for (const navPage of navPages) {

        // заполнение горизонтального меню навигации
        let idFull = "full-nav-link-" + navPage.anchor;
        let navFullPageLine = buildHTML({
            tag: "li",
            class: "hide-on-med-and-down",
            content: tagA(navPage.title, {id: idFull, href: "#nav:" + idFull}),
        });
        $("#navigation-full").append(navFullPageLine);
        $("#"+idFull).bind("click", navPage.onClick);

        // заполнение меню навигации в слайд-панели
        // для добавления горизонтального разделителя между группами меню
        // можно добавить строку типа <li><div class="divider"></div></li>
        let idSlide = "slide-nav-link-" + navPage.anchor;
        let navSlidePageLine = buildHTML({
            tag: "li",
            content: tagA(navPage.title, {id: idSlide, href: "#nav:" + idFull}),
        });
        $("#nav-slide-out").append(navSlidePageLine);
        $("#"+idSlide).bind("click", navPage.onClick);
    }
}