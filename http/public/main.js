window.addEventListener('load', () => {
    document.getElementById('login-form').addEventListener('submit',handleAuth);
    document.getElementById('join-form').addEventListener('submit',handleAuth);
    document.getElementById('post-form').addEventListener('submit',postContent);

    let jwt = Cookie.getCookie('_jwt');
    if (jwt) {
        loadContent();
    }
    else {
        show('login-block');
    }
})

function show(id) {
    const blocks = document.querySelectorAll('.block');
    blocks.forEach((block) =>
        block.classList.add('hidden')
    )
    document.getElementById(id).classList.remove('hidden')
}

function getErrorBlock(form) {
    var b = Array.from(form.children).filter(e => 
        (e.classList.contains('error-block')))
    return b[0];
}

function handleAuth(e) {
    const data = new URLSearchParams(new FormData(e.target)); 

    e.preventDefault();
    getErrorBlock(e.target).innerHTML = ''

    fetch(e.target.action,{method:'post',body: data})
        .then(response => response.json())
        .then(json => handleAuthResponse(e.target, json))
        .catch(error => console.error(error));

}

function handleAuthResponse(form, json) {
    if (json.Errors != null) {
        const e = getErrorBlock(form)
        const ul = document.createElement('ul');
        for (const err of json.Errors) {
            const li = document.createElement('li');
            li.innerText = err;
            ul.append(li);
        }
        e.append(ul);
    }
    else {
        console.log(json);
        loadContent();
    }
}

function loadContent() {
    let jwt = Cookie.getCookie('_jwt');

    fetch('/content',{
        method: "GET",
        headers: {
            'Authorization': 'Bearer ' + jwt 
        }
    })
    .then(response => response.json())
    .then(json => handleContentResponse(json))
    .catch(error => {console.error(error); show('login-block');});
}

function handleContentResponse(json) {
    //console.log('handleContentResponse',json);
    $content = document.getElementById('content-block');
    $content.innerHTML = `<a href="javascript:createPost()">New Post</a>`;
    for (const c of json.Data) {
        $content.innerHTML += 
        `<div class="post">
            <h2>${c.Title} <i>${c.Status == "1" ? "DRAFT" : ""}</i></h2>
            <div>${c.Body}</div>
            <div><i>${c.UserName} ${new Date(c.Created).toDateString()}</i></div>` +
            (c.Access ? `<div><a href="javascript:editPost(${c.ContentId})">Edit</a></div>` : '') +
        '</div>'; 
    }

    show('content-block');
}

function createPost() {
    document.getElementById('id').value = '0'
    document.getElementById('submit-post').value = 'Create Post'
    document.getElementById('post-h2').innerText = 'Create Post'
    document.getElementById('title').value = ''
    document.getElementById('body').value = ''
    document.getElementById('pub-actions').value = ''
    show('post-block')
}

function editPost(id) {
    console.log('editPost',id)
    let jwt = Cookie.getCookie('_jwt');
    let params = new URLSearchParams();
    params.append("id",id.toString())
    fetch('/content',{
            method:'post',
            body: params,
            headers: {
                'Authorization': 'Bearer ' + jwt 
            }})
        .then(response => response.json())
        .then(json => {
            document.getElementById('id').value = id
            document.getElementById('submit-post').value = 'Edit Post'
            document.getElementById('post-h2').innerText = 'Edit Post'
            document.getElementById('title').value = json.Data.Title
            document.getElementById('body').value = json.Data.Body
            if (json.Data.Status == "1") {
                document.getElementById('pub-actions').innerHTML = 
                `<a href="/publish?id=${id}&status=2">Publish</a> <a href="/publish?id=${id}&status=0" class="content-action">Delete</a>`;
            }
            else {
                document.getElementById('pub-actions').innerHTML = 
                `<a href="/publish?id=${id}&status=1">Unpublish</a> <a href="/publish?id=${id}&status=0" class="content-action">Delete</a>`;
            }
            show('post-block')
        })
        .catch(error => console.error(error));

}

function postContent(e) {
    //const $form = document.getElementById(e.target.id);
    const data = new URLSearchParams(new FormData(e.target)); 
    e.preventDefault();
    let jwt = Cookie.getCookie('_jwt');
    //console.log(jwt)
    fetch('/post',{
        method: "post",
        body: data,
        headers: {
            'Authorization': 'Bearer ' + jwt 
        }
    })
    .then(response => response.json())
    .then(json => handleContentResponse(json))
    .catch(error => console.error(error));
}

function logout() {
    document.cookie = '_jwt=;expires=Thu, 01 Jan 1970 00:00:01 GMT';
    location = '/';
}