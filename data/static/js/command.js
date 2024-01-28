function processResponse(resp){
    if(resp.ReturnCode == 200) {
        switch (resp.Command){
            case 'redirect':
                window.location.href=resp.Payload.Target;
                break;
            case 'go_back':
                history.back();
                break;
            //TODO other commands
        }
    }
}