$(document).ready(function () {
    $('table#students').DataTable({
        processing: true,
        serverSide: true,
        ordering: false,
        searching: false,
        ajax: {
            url: config.server + '/v1/users',
            type: 'GET',
            dataSrc: function (d) {
                d.list.forEach(function (element, index, array) {
                    array[index].CreatedAt = element.CreatedAt.substring(0, 19)
                })
                return d.list
            },
            beforeSend: function (xhr) {
                let token = $.cookie('token')
                if (token == undefined) {
                    renewToken()
                    token = $.cookie('token')
                }

                xhr.setRequestHeader('Authorization', 'Bearer ' + token);
            },
            error: function (xhr, error, thrown) {
                if (xhr.status == 401) {
                    let cookies = $.cookie()
                    for (var cookie in cookies) {
                        $.removeCookie(cookie)
                    }

                    location.href = '/login.html'
                } else {
                    alert(xhr.responseText)
                }
            }
        },
        columns: [
            { data: "Name" },
            { data: "Account" },
            { data: "Major" },
            { data: "Number" },
            { data: "CreatedAt" },
        ],
        language: {
            url: '/assets/languages/chinese.json'
        },
    });
})

$('input#upload').fileinput({
    language: 'zh-TW',
    theme: "fas",
    allowedFileExtensions: ['csv'],
    uploadUrl: config.server + '/v1/users',
    ajaxSettings: {
        headers: {
            'Authorization': 'Bearer ' + $.cookie('token'),
        }
    },
}).on('fileuploaderror', function (event, data, msg) {
    $('div.kv-upload-progress.kv-hidden').css({ 'display': 'none' })
})