const TYPE = {
    'both': '同時認列教育實習服務暨志工服務',
    'internship': '實習服務',
    'volunteer': '志工服務',
}
const STATUS = {
    '': '未審核',
    'pass': '通過',
    'failed': '未通過',
}
let StudentServiceLearningID = undefined

const getStudentServiceLearning = () => {
    $.ajax({
        url: `${config.server}/v1/service-learning/sign-up`,
        type: 'GET',
        error: (xhr) => {
            console.error(xhr);
        },
        beforeSend: (xhr) => {
            let token = $.cookie('token')
            if (token == undefined) {
                renewToken()
                token = $.cookie('token')
            }

            xhr.setRequestHeader('Authorization', `Bearer ${token}`);
        },
        success: (response) => {
            if (response.list.length == 0) {
                $('#student-service-learning tbody').append(`
                        <tr>
                            <td scope="row" colspan="8" style="text-align: center">尚無資料</td>
                        </tr>
                    `)
            } else {
                response.list.forEach((element, index) => {
                    let date = element.ServiceLearning.Start.substring(0, 10) + ' ~ ' + element.ServiceLearning.End.substring(0, 10)

                    $('#student-service-learning tbody').append(`
                        <tr>
                            <th scope="row">${index}</th>\
                            <td>${STATUS[element.Status]}</td>\
                            <td>${TYPE[element.ServiceLearning.Type]}</td>\
                            <td>${element.ServiceLearning.Content}</td>\
                            <td>${date}</td>\
                            <td>${element.ServiceLearning.Session}</td>\
                            <td>${element.ServiceLearning.Hours}</td>\
                            <td><button class="btn btn-primary" onclick="edit(${element.ID})">編輯</button></td>\
                        </tr>\
                    `)
                })
            }
        }
    });
}

$(document).ready(() => {
    getStudentServiceLearning()


    $("#reference").fileinput({
        language: 'zh-TW',
        theme: "fas",
        uploadUrl: `${config.server}/v1/service-learning`,
        ajaxSettings: {
            headers: {
                'Authorization': `Bearer ${$.cookie('token')}`,
            },
            method: "PATCH"
        },
        uploadExtraData: (previewId, index) => {
            return {
                'StudentServiceLearningID': StudentServiceLearningID,
            }
        }
    }).on('fileuploaded', (event, previewId, index, fileId) => {
        swal({
            title: '',
            text: '成功',
            icon: "success",
            timer: 1000,
            buttons: false,
        })
    }).on('fileuploaderror', (event, data, msg) => {
        swal({
            title: '',
            text: '失敗',
            icon: "error",
            timer: 1000,
            buttons: false,
        })
    })

    $("#review").fileinput({
        language: 'zh-TW',
        theme: "fas",
        uploadUrl: `${config.server}/v1/service-learning`,
        ajaxSettings: {
            headers: {
                'Authorization': `Bearer ${$.cookie('token')}`,
            },
            method: "PATCH"
        },
        uploadExtraData: (previewId, index) => {
            return {
                'StudentServiceLearningID': StudentServiceLearningID,
            }
        }
    }).on('fileuploaded', (event, previewId, index, fileId) => {
        swal({
            title: '',
            text: '成功',
            icon: "success",
            timer: 1000,
            buttons: false,
        })
    }).on('fileuploaderror', (event, data, msg) => {
        swal({
            title: '',
            text: '失敗',
            icon: "error",
            timer: 1000,
            buttons: false,
        })
    })
})

const edit = id => {
    StudentServiceLearningID = id
    $('#Modal').modal('show')
}