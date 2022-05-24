function getIdmClientDetails() {
    return {
        "id": "{{IDM_CLIENT_ID}}",
        "secret": "{{IDM_CLIENT_SECRET}}",
        "endpoint": "http://am/am/oauth2/realms/root/realms/alpha/access_token",
        "scope": "fr:idm:*",
        "idmAdminUsername": "{{IG_IDM_USER}}",
        "idmAdminPassword": "{{IG_IDM_PASSWORD}}"
    }
}

// Constants
STATUS_AUTHORISED = "Authorised"
var script_name = "policy_evaluation_script.js"
logger.message(script_name + ": starting")

var accountsAndTransactionsPermissions = [{
    name: "READACCOUNTSBASIC",
    property: {permission: "ReadAccountsBasic", requestType: "accounts"}
},
    {name: "READACCOUNTSDETAIL", property: {permission: "ReadAccountsDetail", requestType: "accounts"}},
    {name: "READBALANCES", property: {permission: "ReadBalances", requestType: "balances"}},
    {name: "READBENEFICIARIESBASIC", property: {permission: "ReadBeneficiariesBasic", requestType: "beneficiaries"}},
    {name: "READBENEFICIARIESDETAIL", property: {permission: "ReadBeneficiariesDetail", requestType: "beneficiaries"}},
    {name: "READDIRECTDEBITS", property: {permission: "ReadDirectDebits", requestType: "direct-debits"}},
    {name: "READOFFERS", property: {permission: "ReadOffers", requestType: "offers"}},
    {name: "READPAN", property: {permission: "ReadPAN", requestType: ""}},
    {name: "READPARTY", property: {permission: "ReadParty", requestType: "party"}},
    {name: "READPARTIES", property: {permission: "ReadParty", requestType: "parties"}},
    {name: "READPARTYPSU", property: {permission: "ReadPartyPSU", requestType: "party"}},
    {name: "READPRODUCT", property: {permission: "ReadProducts", requestType: "product"}},
    {name: "READPRODUCTS", property: {permission: "ReadProducts", requestType: "products"}},
    {
        name: "READSCHEDULEDPAYMENTSBASIC",
        property: {permission: "ReadScheduledPaymentsBasic", requestType: "scheduled-payments"}
    },
    {
        name: "READSCHEDULEDPAYMENTSDETAIL",
        property: {permission: "ReadScheduledPaymentsDetail", requestType: "scheduled-payments"}
    },
    {
        name: "READSTANDINGORDERSBASIC",
        property: {permission: "ReadStandingOrdersBasic", requestType: "standing-orders"}
    },
    {
        name: "READSTANDINGORDERSDETAIL",
        property: {permission: "ReadStandingOrdersDetail", requestType: "standing-orders"}
    },
    {name: "READSTATEMENTSBASIC", property: {permission: "ReadStatementsBasic", requestType: "statements"}},
    {name: "READSTATEMENTSDETAIL", property: {permission: "ReadStatementsDetail", requestType: "statements"}},
    {name: "READTRANSACTIONSBASIC", property: {permission: "ReadTransactionsBasic", requestType: "transactions"}},
    {name: "READTRANSACTIONSCREDITS", property: {permission: "ReadTransactionsCredits", requestType: "transactions"}},
    {name: "READTRANSACTIONSDEBITS", property: {permission: "ReadTransactionsDebits", requestType: "transactions"}},
    {name: "READTRANSACTIONSDETAIL", property: {permission: "ReadTransactionsDetail", requestType: "transactions"}}
];

function getPermissionAccountAndTransactions(name) {
    for (var i = 0; i < accountsAndTransactionsPermissions.length; i++) {
        if (accountsAndTransactionsPermissions[i].name === name || accountsAndTransactionsPermissions[i].property.permission == name) {
            return accountsAndTransactionsPermissions[i].property.permission
        }
    }
    return null
}

function getRequestTypeAccountAndTransactions(requestType) {
    for (var i = 0; i < accountsAndTransactionsPermissions.length; i++) {
        if (accountsAndTransactionsPermissions[i].property.requestType === requestType) return accountsAndTransactionsPermissions[i].property.requestType
    }
    return null
}

function dataAuthorisedAccountAndTransactions(permissions, requestType) {
    switch (requestType) {
        case "transactions":
            return ((permissions.indexOf(getPermissionAccountAndTransactions("READTRANSACTIONSBASIC")) > -1 || permissions.indexOf(getPermissionAccountAndTransactions("READTRANSACTIONSDETAIL")) > -1) && (permissions.indexOf(getPermissionAccountAndTransactions("READTRANSACTIONSDEBITS")) > -1 || permissions.indexOf(getPermissionAccountAndTransactions("READTRANSACTIONSCREDITS")) > -1))
        default:
            for (var i = 0; i < accountsAndTransactionsPermissions.length; i++) {
                if (requestType == accountsAndTransactionsPermissions[i].property.requestType && permissions.indexOf(accountsAndTransactionsPermissions[i].property.permission) > -1) return true
            }
    }
    return false
}

function dataAuthorised(permissions, requestType) {
    if (dataAuthorisedAccountAndTransactions(permissions, requestType))
        return true
    return false
}

function parseResourceUri() {
    var elements = resourceURI.split("/");
    return {
        "api": elements[6].indexOf("?") > -1 ? elements[6].substring(0, elements[6].indexOf("?")) : elements[6],
        "id": (elements.length > 7) ? elements[7] : null,
        "data": (elements.length > 8) ? elements[8] : null
    }
}

var p = "=";
var tab = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

function base64encode(ba) {
    var s = [], l = ba.length;
    var rm = l % 3;
    var x = l - rm;
    for (var i = 0; i < x;) {
        var t = ba[i++] << 16 | ba[i++] << 8 | ba[i++];
        s.push(tab.charAt((t >>> 18) & 0x3f));
        s.push(tab.charAt((t >>> 12) & 0x3f));
        s.push(tab.charAt((t >>> 6) & 0x3f));
        s.push(tab.charAt(t & 0x3f));
    }
    //	deal with trailers, based on patch from Peter Wood.
    switch (rm) {
        case 2: {
            var t = ba[i++] << 16 | ba[i++] << 8;
            s.push(tab.charAt((t >>> 18) & 0x3f));
            s.push(tab.charAt((t >>> 12) & 0x3f));
            s.push(tab.charAt((t >>> 6) & 0x3f));
            s.push(p);
            break;
        }
        case 1: {
            var t = ba[i++] << 16;
            s.push(tab.charAt((t >>> 18) & 0x3f));
            s.push(tab.charAt((t >>> 12) & 0x3f));
            s.push(p);
            s.push(p);
            break;
        }
    }
    return s.join("");	//	string
}

function base64decode(str) {
    var s = str.split(""), out = [];
    var l = s.length;
    while (s[--l] == p) {
    }	//	strip off trailing padding
    for (var i = 0; i < l;) {
        var t = tab.indexOf(s[i++]) << 18;
        if (i <= l) {
            t |= tab.indexOf(s[i++]) << 12
        }
        ;
        if (i <= l) {
            t |= tab.indexOf(s[i++]) << 6
        }
        ;
        if (i <= l) {
            t |= tab.indexOf(s[i++])
        }
        ;
        out.push((t >>> 16) & 0xff);
        out.push((t >>> 8) & 0xff);
        out.push(t & 0xff);
    }
    //	strip off any null bytes
    while (out[out.length - 1] == 0) {
        out.pop();
    }
    return out;	//	byte[]
}

function stringFromArray(data) {
    var count = data.length;
    var str = "";

    for (var index = 0; index < count; index += 1)
        str += String.fromCharCode(data[index]);

    return str;
}

function logResponse(callerMethod, response) {
    logger.message(script_name + ": [" + callerMethod + "] OB_Policy User REST Call. Status: " + response.getStatus() + ", Body: " + response.getEntity().getString());
}

function getIdmAccessToken() {

    var clientInfo = getIdmClientDetails();
    var request = new org.forgerock.http.protocol.Request();
    request.setUri(clientInfo.endpoint);
    request.setMethod("POST");
    request.getHeaders().add("Content-Type", "application/x-www-form-urlencoded");
    var formvars = "grant_type=password" +
        "&client_id=" + clientInfo.id +
        "&client_secret=" + clientInfo.secret +
        "&scope=" + clientInfo.scope +
        "&username=" + clientInfo.idmAdminUsername +
        "&password=" + clientInfo.idmAdminPassword;
    request.setEntity(formvars);

    var response = httpClient.send(request).get();


    logResponse("getIdmAccessToken", response);

    var oauth2response = JSON.parse(response.getEntity().getString());

    var accessToken = oauth2response.access_token
    logger.message(script_name + ": Got access token " + accessToken);
    return accessToken
}

function findIntentType(api) {
    if (getRequestTypeAccountAndTransactions(api) != null)
        return "accountAccessIntent"
    else if (api === "domestic-payments")
        return "domesticPaymentIntent"
    return null
}

function getIntent(intentId, intentType) {
    var accessToken = getIdmAccessToken();
    var request = new org.forgerock.http.protocol.Request();
    var uri = "http://idm/openidm/managed/" + intentType + "/" + intentId + "?_fields=_id,_rev,Data,Risk,user/_id,accounts,apiClient/_id"
    logger.message(script_name + ": IDM fetch " + uri)

    request.setMethod('GET');
    request.setUri(uri)
    request.getHeaders().add("Authorization", "Bearer " + accessToken);

    var response = httpClient.send(request).get();
    logResponse("getIntent", response);

    var intent = JSON.parse(response.getEntity().getString());
    return intent
}


function initiationMatch(initiationRequest, initiation) {

    // TODO: do comparison at object level, like JSONAssert()

    var initiationRequestObj = JSON.parse(stringFromArray(base64decode(initiationRequest)))

    var match =
        (initiationRequestObj.InstructionIdentification == initiation.InstructionIdentification) &&
        (initiationRequestObj.EndToEndIdentification == initiation.EndToEndIdentification) &&
        (initiationRequestObj.InstructedAmount.Amount == initiation.InstructedAmount.Amount) &&
        (initiationRequestObj.InstructedAmount.Currency == initiation.InstructedAmount.Currency) &&
        (initiationRequestObj.CreditorAccount.SchemeName == initiation.CreditorAccount.SchemeName) &&
        (initiationRequestObj.CreditorAccount.Name == initiation.CreditorAccount.Name) &&
        (initiationRequestObj.CreditorAccount.SecondaryIdentification == initiation.CreditorAccount.SecondaryIdentification) &&
        (initiationRequestObj.RemittanceInformation.Reference == initiation.RemittanceInformation.Reference) &&
        (initiationRequestObj.RemittanceInformation.Unstructured == initiation.RemittanceInformation.Unstructured)


    if (!match) {
        logger.warning(script_name + ": Mismatch between request [" + JSON.stringify(initiationRequestObj) + "] and consent [" + JSON.stringify(initiation) + "]");
    }

    return match
}

var intentId = environment.get("intent_id").iterator().next();
var apiRequest = parseResourceUri()
logger.message(script_name + ": req " + apiRequest.api + ":" + apiRequest.id + ":" + apiRequest.data);

var intentType = findIntentType(apiRequest.api)
var intent = getIntent(intentId, intentType);

if (intentType === "accountAccessIntent") {
    logger.message(script_name + ": Account Access Intent");

    var status = intent.Data.Status
    var permissions = intent.Data.Permissions
    var accounts = intent.accounts
    // The responseAttributes expected always and array as value
    var userResourceOwner = new Array(intent.user._id)

    if (status != STATUS_AUTHORISED) {
        logger.message(script_name + "-[Account Access]: Rejecting request - status [" + status + "]")
        authorized = false
    } else if (apiRequest.id == null) {
        logger.message(script_name + "-[Account Access]: accounts " + accounts);
        responseAttributes.put("grantedAccounts", accounts);
        responseAttributes.put("grantedPermissions", permissions);
        responseAttributes.put("userResourceOwner", userResourceOwner);
        authorized = true
    } else if (apiRequest.data == null) {
        logger.message(script_name + "-[Account Access]: account info for " + apiRequest.id);
        // RS server expects granted accounts and permissions even though we're checking as well
        responseAttributes.put("grantedAccounts", accounts);
        responseAttributes.put("grantedPermissions", permissions);
        responseAttributes.put("userResourceOwner", userResourceOwner);
        authorized = (accounts.indexOf(apiRequest.id) > -1) &&
            dataAuthorised(permissions, apiRequest.api)
    } else {
        logger.message(script_name + "-[Account Access]: account request for " + apiRequest.id + ":" + apiRequest.data);
        // RS server expects granted accounts and permissions even though we're checking as well
        responseAttributes.put("grantedAccounts", accounts);
        responseAttributes.put("grantedPermissions", permissions);
        responseAttributes.put("userResourceOwner", userResourceOwner);
        authorized = (accounts.indexOf(apiRequest.id) > -1) &&
            dataAuthorised(permissions, apiRequest.data)
    }

} else if (intentType === "domesticPaymentIntent") {
    logger.message(script_name + ": Domestic Payments Intent");

    var status = intent.Data.Status
    var userResourceOwner = new Array(intent.user._id)

    if (status !== STATUS_AUTHORISED) {
        logger.message(script_name + "-[Domestic Payments]: Rejecting request - status [" + status + "]")
        authorized = false
    } else {
        responseAttributes.put("userResourceOwner", userResourceOwner);
        var requestMethod = environment.get("request_method").iterator().next()

        if (requestMethod != null) {
            switch (requestMethod) {
                case "POST":
                    logger.message(script_name + "-[Domestic Payments]: POST request");
                    var initiation = environment.get("initiation").iterator().next()
                    authorized = initiationMatch(initiation, intent.Data.Initiation)
                    break;
                case "GET":
                    logger.message(script_name + "-[Domestic Payments]: GET request");
                    authorized = true
                    break;
                default:
                    authorized = false
                    break;
            }
        } else {
            authorized = false
        }
    }
} else {
    authorized = false
}
logger.message(script_name + ": Policy evaluation result, authorized=" + authorized);