
let tabs = document.getElementsByClassName("tab");
for (let i = 1; i < tabs.length; i++) {
    tabs[i].style.display="none";
}
// 切换页面
for (let item of document.getElementsByClassName("nav-link")) {
    let idx = item.getAttribute("index");
    item.addEventListener("click", function() {
        for (let i = 0; i < tabs.length; i++) {
            if ((""+i) === idx) {
                tabs[i].style.display = "block";
            } else {
                tabs[i].style.display = "none";
            }
        }
    });
}

function sendFile(ip){
    // alert(ip)
    //创建文件输入对象
    let fileInput = document.createElement("input");
    fileInput.type = "file";
    fileInput.multiple = true;
    fileInput.setAttribute('multiple', '')
    fileInput.onchange = function () {
        let files = fileInput.files;
        let formData = new FormData();  // 创建FormData对象

        formData.append("ip", ip);
        for (let i = 0; i < files.length; i++) {
            formData.append('files[]', files[i]);
        }

        // 使用fetch API发送文件
        fetch('/upload', {
            method: 'POST',  // 设置请求方法为POST
            body: formData,  // 设置请求主体为包含文件的FormData对象
        })
        .then(response => {
            if (response.ok) {
                console.log('文件上传成功！');
                toast("文件上传成功");
                setTimeout(function (){
                    window.location = "/"
                }, 2000)
            } else {
                console.error('文件上传失败：' + response.statusText);
                toast("文件上传失败");
            }
        })
        .catch(error => {
            console.error('文件上传失败：' + error);
            toast("文件上传失败");
        });
    }
    fileInput.click();
}


function rename(){
    let name = document.getElementById("nick-name").value;
    fetch("/rename/" + name)
        .then(res => {
            if (!res.ok) {
                throw new Error('Network response was not ok');
            }
            return res.json(); // 将响应体解析为JSON格式
        }).then(res => {
            if (res.code === 0){
                document.getElementById("show-nick-name").innerText = name
                toast("修改成功");
            }else {
                toast("修改失败");
            }
        })
        .catch(err => {
            console.log(err)
            toast("修改失败");
        })
}


function setMyInfo(){
    fetch("/myHost")
        .then(res => {
            if (!res.ok) {
                throw new Error('Network response was not ok');
            }
            return res.json(); // 将响应体解析为JSON格式
        }).then(res => {
            if (res.code === 0){
                document.getElementById("show-nick-name").innerText = res.data.nickName;

            }else {

            }
        })
        .catch(err => {
            console.log(err)
        })
}

setMyInfo()

function setDestIp(ip){
    document.getElementById("destIp").innerText = ip
}

function sendMessage(){
    let formData = new FormData();  // 创建FormData对象
    // alert(document.getElementById("message").value)
    // alert(document.getElementById("destIp").innerText)
    formData.append("msg", document.getElementById("message").value);
    formData.append("ip",  document.getElementById("destIp").innerText);
    // 使用fetch API发送文件
    fetch('/message', {
        method: 'POST',  // 设置请求方法为POST
        body: formData,  // 设置请求主体为包含文件的FormData对象
    })
    .then(response => {
        if (response.ok) {
            console.log('消息发送成功！');
            toast("消息发送成功");
            setTimeout(function (){
                window.location = "/"
            }, 2000)
        } else {
            console.error('消息发送失败：' + response.statusText);
            toast("消息发送失败");
        }
    })
    .catch(error => {
        console.error('消息发送失败：' + error);
        toast("消息发送失败");
    });
}

function deleteAll(id){
    fetch("/deleteFile/" + id)
        .then(res => {
            if (!res.ok) {
                throw new Error('Network response was not ok');
            }
            return res.json(); // 将响应体解析为JSON格式
        }).then(res => {
        if (res.code === 0){
            console.log("ok")
            toast("删除成功");
            setTimeout(function (){
                window.location = "/"
            }, 2000)
        }else {

            toast("删除失败");
        }
    })
        .catch(err => {
            toast("删除失败");
            console.log(err)
        })
}

function download(id, name){
    // alert(id)
    let url = "/downloadFile/" + id + "/" + name

    // window.open("/downloadFile/" + id + "/" + e);

    // 创建一个新的 iframe 元素
    let iframe = document.createElement('iframe');

    // 将 iframe 的 'src' 属性设置为文件的 URL
    iframe.src = url;

    // 将 iframe 设置为隐藏
    iframe.style.display = 'none';

    // 将 iframe 添加到页面中
    document.body.appendChild(iframe);
    // 一段时间后移除这些 iframe
    setTimeout(() => {
        document.removeChild(iframe)
    }, 5000);
}

function downloadAll(id){
    // alert(id)
    let url = "/download/" + id

    // window.open("/downloadFile/" + id + "/" + e);

    // 创建一个新的 iframe 元素
    let iframe = document.createElement('iframe');

    // 将 iframe 的 'src' 属性设置为文件的 URL
    iframe.src = url;

    // 将 iframe 设置为隐藏
    iframe.style.display = 'none';

    // 将 iframe 添加到页面中
    document.body.appendChild(iframe);
    // 一段时间后移除这些 iframe
    setTimeout(() => {
        document.removeChild(iframe)
    }, 5000);
}

function copyToClipboard(msg) {


    var aux = document.createElement("input");
    aux.setAttribute("value", msg);
    document.body.appendChild(aux);
    aux.select();
    document.execCommand("copy");
    document.body.removeChild(aux);
    toast("已将内容复制到黏贴版");
}

function toast(msg){
    let t = document.getElementById("toast");
    t.innerText = msg;
    const toastLiveExample = document.getElementById('liveToast');
    const toast = new bootstrap.Toast(toastLiveExample);
    toast.show();
}