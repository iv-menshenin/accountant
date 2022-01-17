
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

function AccountsList() {
    let showFn = preparePage("Лицевые счета", [
        {
            tag: "div", class: "", content: [
                {id: "accounts-container", tag: "ul", class: ["collection", "with-header"], content: {tag: "li", class: "collection-header", content: {tag: "h4", content: "Зарегистрированные ЛС"}}},
            ]
        },
        {
            tag: "button", "data-target": "modal-action", class: ["btn", "waves-effect", "waves-light", "modal-trigger", "action-button"], content: "Добавить новый"
        }
    ], ()=>{
        let accountsCollection = accounts.getAccounts();
        $("#accounts-container").append(accountsCollection.map(renderAccountListElement));
    });
    prepareModalForm([
        {
            id: "account-add-form", tag: "div", class: ["container"], content: [
                {tag: "h5", content: "Новый лицевой счет"},
                {
                    tag: "div", class: "row", content: [
                        {
                            tag: "div", class: ["col", "s10"], content: [
                                {tag: "label", for: "account-number", content: "Лицевой счет"},
                                {id: "account-number", tag: "input", type: "text"},
                            ]
                        }
                    ]
                },
            ]
        }
    ]);
    accounts.loadAccounts(
        (accounts) => {
            showFn();
        },
        (error) => {
            // todo toss error
        },
    )
}

function showAccountDetails(anchor) {
    let account_id = undefined;
    let url = new URL(anchor.href);
    let uuid_match = /#account:([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})/
    let hash_blocks = uuid_match.exec(url.hash);
    if (hash_blocks.length > 1) {
        account_id = hash_blocks[1]
    }
    if (account_id) {
        let account = accounts.getAccount(account_id);
        if (account) {

            return true
        }
    }
    // todo toss error
}