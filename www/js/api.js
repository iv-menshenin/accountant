
class apiManager {
    session = "";
    server = "";
    port = "";
    proto = "https";
    rooPath = "/api";
    requestedWith = "X-SessionManager";

    constructor() {
        if (document.domain === "localhost") {
            this.server = "localhost";
            this.port = "8080";
            this.proto = "http";
        } else {
            this.server = document.domain;
            this.port = "";
            this.proto = "https";
        }
        this.session = localStorage.getItem("session")
        if (!this.session) {
            // todo login
        }
    }

    getEncodedUriParameterString(paramsMap){
        if (paramsMap) {
            let keys = Object.keys(paramsMap);
            let outResult = "";
            for (let nn = 0;nn < keys.length;nn++) {
                outResult += (outResult ? "&" : "") + keys[nn] + "=" + encodeURIComponent(paramsMap[keys[nn]]);
            }
            return outResult;
        } else {
            return ""
        }
    }

    getApiAddr(path, query) {
        let encodedQuery = "";
        if (query) {
            encodedQuery = this.getEncodedUriParameterString(query);
        }
        let server = this.server;
        if (this.port) {
            server = server + ":" + this.port;
        }
        return this.proto + "://" + server + this.rooPath + path + "?" + encodedQuery
    }

    apiExecute(path, method, query, data, onSuccess, onError) {
        let self = this;
        let headers = {
            "Accept": "application/json",
            "X-Requested-With": this.requestedWith,
        };
        $.ajax({
            crossDomain: true,
            type: method,
            url: self.getApiAddr(path, query),
            headers: headers,
            data: (method !== "GET" ? JSON.stringify(data) : undefined),
            success: function (responseData) {
                onSuccess(responseData)
            },
            error: function (request, _, cl) {
                onError(request)
                console.log(cl);
                if (request.status === 401) {
                    // todo login page redirect
                }
            }
        });
    }

    // API -------------------------------------------------------------------------------------------------------------

    accountsPath = "/accounts";

    // GetAccounts получает список зарегистрированых ЛС
    GetAccounts(onSuccess, onError) {
        let query = {}
        this.apiExecute(this.accountsPath, "GET", query, undefined, (data) => {
            onSuccess(data.data)
        }, onError);
    }

    // CreateAccount создает новый лицевой счет в ответе на запрос присылает полную структуру данных созданного ЛС
    CreateAccount(account, onSuccess, onError) {
        let body = {
            account: account.account,
            cad_number: account.cad_number,
            agreement: account.agreement,
            agreement_date: account.agreement_date,
            comment: account.comment,
        };
        this.apiExecute(this.accountsPath, "POST", undefined, body, onSuccess, onError);
    }

    // UpdateAccount обновляет ЛС указанными полями
    UpdateAccount(account, onSuccess, onError) {
        let path = this.accountsPath + "/" + account.account_id;
        let body = {
            account: account.account,
            cad_number: account.cad_number,
            agreement: account.agreement,
            agreement_date: account.agreement_date,
            comment: account.comment,
        };
        this.apiExecute(path, "PUT", undefined, body, onSuccess, onError);
    }

    // DeleteAccount удаляет ЛС по его идентификатору
    DeleteAccount(account_id, onSuccess, onError) {
        let path = this.accountsPath + "/" + account_id;
        this.apiExecute(path, "DELETE", undefined, undefined, onSuccess, onError);
    }
}

let manager = new(apiManager);
