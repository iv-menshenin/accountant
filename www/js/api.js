
class apiManager {
    server = "";
    port = "";
    proto = "https";
    rooPath = "/api";
    requestedWith = "X-SessionManager";

    authHeader = "X-Auth-Token";
    authToken  = "";
    refreshToken  = "";
    context = [];

    constructor() {
        if (document.domain === "localhost") {
            this.server = "localhost";
            this.port = "8080";
            this.proto = "http:";
        } else {
            this.server = document.location.hostname;
            this.port = document.location.port;
            this.proto = document.location.protocol;
        }
        this.refreshToken = localStorage.getItem('rt');
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
            encodedQuery = "?" + this.getEncodedUriParameterString(query);
        }
        let server = this.server;
        if (this.port) {
            server = server + ":" + this.port;
        }
        return this.proto + "//" + server + this.rooPath + path + encodedQuery
    }

    apiExecute(path, method, query, data, onSuccess, onError) {
        let self = this;
        let headers = {
            "Accept": "application/json",
            "X-Requested-With": this.requestedWith,
        };
        let exec = ()=>{};
        let suc = function (responseData) {
            onSuccess(responseData.data)
        };
        let err = function (req, cl, canRepeat) {
            if (req.status === 401 && canRepeat) {
                self.AuthInitiate(exec);
                return
            }
            if (req.responseJSON) {
                console.log(req.responseJSON.meta);
                onError(req.responseJSON.meta.message, req.status)
                return
            }
            if (req.responseText) {
                console.log(req.responseText);
                onError(req.responseText, req.status)
                return
            }
            onError(req, req.status)
            console.log(cl);
        };
        exec = (canRepeat)=> {
            headers[this.authHeader] = this.authToken;
            $.ajax({
                crossDomain: true,
                type: method,
                url: self.getApiAddr(path, query),
                headers: headers,
                data: (method !== "GET" ? JSON.stringify(data) : undefined),
                success: suc,
                error: (req, _, cl) => err(req, cl, canRepeat)
            });
        }
        exec(true);
    }

    // AUTH ------------------------------------------------------------------------------------------------------------

    AuthInitiate(callback) {
        localStorage.removeItem('rt');
        if (this.refreshToken) {
            this.Refresh(()=>{
                callback(false);
            }, ()=>{
                document.location.replace("#login:back=" + document.location.hash.substr(1));
            })
        } else {
            document.location.replace("#login:back=" + document.location.hash.substr(1));
        }
    }

    loginPath = "/auth/login";

    // Login аутентифицирует пользователя по логину и паролю
    Login(login, password, onSuccess, onError) {
        let body = {
            login: login,
            password: password,
        };
        let self = this;
        let headers = {
            "Accept": "application/json",
            "X-Requested-With": this.requestedWith,
        };
        $.ajax({
            crossDomain: true,
            type: "POST",
            url: this.loginPath,
            headers: headers,
            data: JSON.stringify(body),
            success: function (responseData) {
                self.authToken = responseData.data.jwt_token;
                self.refreshToken = responseData.data.refresh_token;
                self.context = responseData.data.context;
                localStorage.setItem('rt', responseData.data.refresh_token);
                onSuccess(responseData.data.user_id)
            },
            error: function (req, _, cl) {
                if (req.status > 399 && req.status < 500) {
                    onError("возможно вы ввели неверный логин или пароль")
                    return
                }
                onError("неожиданная ошибка")
            }
        });
    }

    refreshPath = "/auth/refresh";

    // Login аутентифицирует пользователя по логину и паролю
    Refresh(onSuccess, onError) {
        let body = {
            token: this.refreshToken,
        };
        let self = this;
        let headers = {
            "Accept": "application/json",
            "X-Requested-With": this.requestedWith,
        };
        $.ajax({
            crossDomain: true,
            type: "POST",
            url: this.refreshPath,
            headers: headers,
            data: JSON.stringify(body),
            success: function (responseData) {
                self.authToken = responseData.data.jwt_token;
                self.refreshToken = responseData.data.refresh_token;
                self.context = responseData.data.context;
                localStorage.setItem('rt', responseData.data.refresh_token);
                onSuccess(responseData.data.user_id)
            },
            error: function (req, _, cl) {
                self.AuthInitiate();
            }
        });
    }

    // API -------------------------------------------------------------------------------------------------------------

    accountsPath = "/accounts";

    // GetAccounts получает список зарегистрированых ЛС
    GetAccounts(onSuccess, onError) {
        this.apiExecute(this.accountsPath, "GET", undefined, undefined, onSuccess, onError);
    }

    // CreateAccount создает новый лицевой счет в ответе на запрос присылает полную структуру данных созданного ЛС
    CreateAccount(account, onSuccess, onError) {
        let body = {
            account: account.account,
            cad_number: account.cad_number,
            agreement: account.agreement,
            agreement_date: account.agreement_date,
            purchase_kind: account.purchase_kind,
            purchase_date: account.purchase_date,
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
            purchase_kind: account.purchase_kind,
            purchase_date: account.purchase_date,
            comment: account.comment,
        };
        this.apiExecute(path, "PUT", undefined, body, onSuccess, onError);
    }

    // DeleteAccount удаляет ЛС по его идентификатору
    DeleteAccount(account_id, onSuccess, onError) {
        let path = this.accountsPath + "/" + account_id;
        this.apiExecute(path, "DELETE", undefined, undefined, onSuccess, onError);
    }

    personsPath = "/persons";

    // CreatePerson создает новый лицевой счет в ответе на запрос присылает полную структуру данных созданного ЛС
    CreatePerson(accountID, person, onSuccess, onError) {
        let body = {
            name: person.name,
            surname: person.surname,
            pat_name: person.pat_name,
            dob: person.dob,
            is_member: person.is_member,
            phone: person.phone,
            email: person.email,
        };
        this.apiExecute(this.accountsPath + "/" + accountID + this.personsPath, "POST", undefined, body, onSuccess, onError);
    }

    // GetAccounts получает список зарегистрированых ЛС
    GetAccountPersons(accountID, onSuccess, onError) {
        this.apiExecute(this.accountsPath + "/" + accountID + this.personsPath, "GET", undefined, undefined, onSuccess, onError);
    }
}

let manager = new(apiManager);
