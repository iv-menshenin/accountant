
class accountsManager {
    collection = [];

    loadAccounts(onDone, onError) {
        manager.GetAccounts(
            (accounts) => {
                this.collection = accounts;
                onDone(accounts);
            },
            (err) => {
                if (err.status === 404) {
                    onDone([]);
                    return
                }
                onError(err.responseJSON.meta.message);
            },
        )
    }

    getAccount(account_id) {
        let filtered = this.collection.filter((account) => {return account.account_id === account_id});
        if (filtered.length === 1) {
            return filtered[0]
        }
        return undefined
    }

    getAccounts() {
        return this.collection;
    }
}

let accounts = new(accountsManager);

function AccountsListPage() {
    let switcher = makeSwitcher();
    let collection = makeCollectionContainer("Зарегистрированные ЛС", {});
    let showFn = preparePage("Лицевые счета", [
        collection.content,
        makeButton("Добавить новый", {target: "modal-action", class: "action-button"}),
    ], ()=>{
        let accountsCollection = accounts.getAccounts();
        collection.append(accountsCollection.map(mapAccountToListElement));
        switcher.switch(true);
    });
    let newAccountModal = accountEditor.pageRender({}, {switch: switcher.consume});
    prepareModalForm([
        {
            id: "account-add-form", tag: "div", class: ["container"], content: [
                windowHeader("Новый лицевой счет"),
                {tag: "div", class: "row", content: newAccountModal.content},
            ]
        }
    ], newAccountModal.footer);
    accounts.loadAccounts(
        (accounts) => {
            showFn();
        },
        (error) => {
            // todo toss error
        },
    )
}

function AccountPage(props, retry=true) {
    if (!props.uuid) {
        return false
    }
    let account = accounts.getAccount(props.uuid);
    if (!account) {
        if (retry) {
            accounts.loadAccounts(()=> {AccountPage(props, false)}, (err) => {});
            return false
        }
        return false
    }
    let accountEditModal = accountEditor.pageRender(account);
    let showFn = preparePage(getFirstPersonName(account), [
        accountEditModal.content,
        {
            tag: "div", class: "row", content: [
                {tag: "button", "data-target": "modal-action-1", class: ["btn", "waves-effect", "waves-light", "modal-trigger", "action-button", "s4", "col"], content: "Собственники"},
                {tag: "div", class: ["col", "s4"]},
                {tag: "button", "data-target": "modal-action-2", class: ["btn", "waves-effect", "waves-light", "modal-trigger", "action-button", "s4", "col"], content: "Участки"}
            ]
        }
    ], ()=>{

    });
    if (account) {
        showFn();
        return true
    }
}