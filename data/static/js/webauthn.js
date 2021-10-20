function step1(){
$.get('/webauthn/' + 'uaerparam', {
   /* attType: attestation_type,
    authType: authenticator_attachment,
    userVerification: user_verification,
    residentKeyRequirement: resident_key_requirement,
    txAuthExtension: txAuthSimple_extension,*/
}, null, 'json')
    .done(function (makeCredentialOptions) {
        makeCredentialOptions.publicKey.challenge = bufferDecode(makeCredentialOptions.publicKey.challenge);
        makeCredentialOptions.publicKey.user.id = bufferDecode(makeCredentialOptions.publicKey.user.id);
        if (makeCredentialOptions.publicKey.excludeCredentials) {
            for (var i = 0; i < makeCredentialOptions.publicKey.excludeCredentials.length; i++) {
                makeCredentialOptions.publicKey.excludeCredentials[i].id = bufferDecode(makeCredentialOptions.publicKey.excludeCredentials[i].id);
            }
        }
        console.log("Credential Creation Options");
        console.log(makeCredentialOptions);
        navigator.credentials.create({
            publicKey: makeCredentialOptions.publicKey
        }).then(function (newCredential) {
            console.log("PublicKeyCredential Created");
            console.log(newCredential);
            state.createResponse = newCredential;
            registerNewCredential(newCredential);
        }).catch(function (err) {
            console.info(err);
        });
    });
}

// Don't drop any blanks
// decode
function bufferDecode(value) {
    return Uint8Array.from(atob(value), c => c.charCodeAt(0));
}