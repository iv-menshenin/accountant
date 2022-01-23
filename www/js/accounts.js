
class accountsManager {
    onDone = () => {};
    collection = [];

    loadAccounts(onDone, onError) {
        this.onDone = onDone;
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

    addOrReplaceAccount(newAccount) {
        let self = this;
        let replaced = this.collection.reduce((replaced, account, i)=>{
            if (account.account_id === newAccount.account_id) {
                self.collection[i] = newAccount;
                return true;
            }
            return replaced;
        }, false)
        if (!replaced) {
            this.collection.push(newAccount);
        }
        this.onDone(this.collection);
    }
}

let accounts = new(accountsManager);

function AccountsListPage() {
    let switcher = makeSwitcher();
    let collection = makeCollectionContainer("Зарегистрированные ЛС", {});
    let showFn = preparePage("Лицевые счета", [
        collection.content,
        makeButton("Добавить новый", {target: "modal-action", class: "action-button"}),
    ], (accounts)=>{
        collection.clear();
        collection.append(accounts.map(mapAccountToListElement));
        switcher.switch(true);
    });
    let newAccountModal = accountEditor.pageRender({}, {switch: switcher.consume, onDone: ()=>{modalWindow.close();}});
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
            showFn(accounts);
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
    let personsContainer = makeCollectionContainer("Собственники", {});
    let objectsContainer = makeCollectionContainer("Объекты", {});
    let showFn = preparePage(getFirstPersonName(account), [
        accountEditModal.content,
        {tag: "div", class: "row", content: [{tag: "hr", class: ["s12", "col"]}]},
        personsContainer.content,
        {
            tag: "div", class: "row", content: [
                {tag: "button", "data-target": "modal-action-1", class: ["btn", "waves-effect", "waves-light", "modal-trigger", "action-button", "s4", "col", "right"], content: "+собственник"}
            ]
        },
        objectsContainer.content,
        {
            tag: "div", class: "row", content: [
                {tag: "button", "data-target": "modal-action-2", class: ["btn", "waves-effect", "waves-light", "modal-trigger", "action-button", "s4", "col", "right"], content: "+участок"}
            ]
        }
    ], ()=>{
        if (account.persons) {
            personsContainer.clear();
            personsContainer.append(account.persons.map(mapPersonToListElement));
        }
        if (account.objects) {
            objectsContainer.clear();
            objectsContainer.append(account.objects.map(mapObjectToListElement));
        }
    });
    if (account) {
        showFn();
        return true
    }
}