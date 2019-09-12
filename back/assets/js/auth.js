$(document).ready(function () {
    let token = $.cookie('token')
    let refreshToken = $.cookie('refresh-token')

    if (token == undefined) {
        if (refreshToken != undefined) {
            renewToken()
            token = $.cookie('token')
        } else if (window.location.pathname != '/login.html') {
            location.href = '/login.html'
        }
    } else if (window.location.pathname == '/login.html') {
        location.href = '/'
    }
})

$('#loginform').on('submit', function (e) {
    e.preventDefault();

    let account = $('input#account').val()
    let password = $('input#password').val()

    $.ajax({
        url: config.server + '/v1/login',
        type: 'POST',
        data: {
            'Account': account,
            'Password': password,
            'Role': 'admin',
        },
        error: function (xhr) {
            alert('Unexcepted Error')
            console.error(xhr);
        },
        success: function (response) {
            if (response.code != 0) {
                alert(response.message)
            } else {
                let date = new Date()
                date.setTime(date.getTime() + (response.data.Expire * 1000));

                $.cookie('token', response.data.Token, {
                    expires: date,
                });
                $.cookie('account', account);
                $.cookie('refresh-token', response.data.RefreshToken)

                location.href = '/'
            }
        }
    });
});

$('button#logout').click(function () {
    $('#logoutModal').show()
})

$('#logoutModal button.btn.btn-primary').click(function () {
    $.ajax({
        url: config.server + '/v1/logout',
        type: 'POST',
        error: function (xhr) {
            console.error(xhr);
        },
        beforeSend: function (xhr) {
            let token = $.cookie('token')
            if (token == undefined) {
                renewToken()
                token = $.cookie('token')
            }

            xhr.setRequestHeader('Authorization', 'Bearer ' + token);
        },
        success: function (response) {
            if (response.code != 0) {
                console.error(response.message)
            } else {
                let date = new Date()
                date.setTime(date.getTime() + (response.data.Expire * 1000));

                let cookies = $.cookie()
                for (var cookie in cookies) {
                    $.removeCookie(cookie)
                }

                location.href = '/login.html'
            }
        }
    });
})

function renewToken() {
    let account = $.cookie('account')
    let refreshToken = $.cookie('refresh-token')

    $.ajax({
        url: config.server + '/v1/renew-token',
        type: 'POST',
        async: false,
        cache: false,
        data: {
            Account: account,
            RefreshToken: refreshToken,
        },
        error: function (xhr) {
            let cookies = $.cookie()
            for (var cookie in cookies) {
                $.removeCookie(cookie)
            }

            location.href = '/login.html'
        },
        success: function (response) {
            if (response.code != 0) {
                let cookies = $.cookie();
                for (var cookie in cookies) {
                    $.removeCookie(cookie)
                }

                location.href = '/login.html'
            } else {
                let date = new Date()
                date.setTime(date.getTime() + (response.data.Expire * 1000));

                $.cookie('token', response.data.Token, {
                    expires: date,
                });
            }
        }
    });
}