//////////////////////////////////////////////////////////////////////////////////////////////////////
//                                          Sessionmanagment                                        //
//////////////////////////////////////////////////////////////////////////////////////////////////////
/**
 * 
 * @param {Event} event
 */
function login(event) {

    try {

        //Get Formdata
        toggleLoadingScreen();
        event.preventDefault();
        var formdata = event.target;

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open (formdata.method, formdata.action, true);
        xhr.send (new FormData(formdata));

        var success = function(response) {

            if(response.status === 200) {

                getUserdata();
                setupLoggedPage();
                dissmissAllModals();

            }
            if(response.status === 409) {

                var messages = "";

                jsonMessages = JSON.parse(response.responseText).Messages;

                for (let index = 0; index < jsonMessages.length; index++) {
                    const element = jsonMessages[index];

                    messages += element +"\n";

                }

                alert(messages);
                setTimeout(function(){dissmissAllModals()}, 1000);

            }
            if(response.status === 500) {

                alert("Something Made wrong: Message under Construction ;)")
                setTimeout(function(){dissmissAllModals()}, 1000);

            }

        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();        
        
        };

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}document.getElementById("loginForm").addEventListener('submit', function(ev){login(ev);})


/**
 * 
 * @param {Event} event 
 */
function register(event) {

    try {

        //Get Formdata
        toggleLoadingScreen();
        event.preventDefault();
        var formdata = event.target;

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open (formdata.method, formdata.action, true);
        xhr.send (new FormData(formdata));

        var success = function(response) {

            if(response.status === 200) {

                getUserdata();
                setupLoggedPage();
                toggleLoadingScreen();
                dissmissAllModals();

            }
            if(response.status === 409) {

                var messages = "";

                jsonMessages = JSON.parse(response.responseText).Messages;

                for (let index = 0; index < jsonMessages.length; index++) {
                    const element = jsonMessages[index];

                    messages += element +"\n";

                }

                alert(messages);
                setTimeout(function(){dissmissAllModals()}, 1000);

            }
            if(response.status === 500) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }

        }
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();
        }

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}document.getElementById("registerForm").addEventListener('submit', function(ev){register(ev);})


/**
 * 
 * @param {HTMLFormElement} formdata 
 */
function logout(event) {

    try {

        //Get Formdata
        event.preventDefault();
        var formdata = event.target;
        formdata.action = `/users/${applicationState.userID}?action=logout`

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open (formdata.method, formdata.action, true);
        xhr.send ();

        //Callbacks
        var success = function(response) {

            setupLoggedOutPage();
            unsetUserdata();

        };
        var error  = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();

        };

    } catch {

        toggleLoadingScreen();

    }

}document.getElementById("logoutForm").addEventListener('submit', function(ev){logout(ev);})


/**
 * 
 */
function getUserdata() {

    try {

        //Get Action
        var action = "/users?action=userdata";
        var method = "GET";

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open (method, action, false);

        var success = function(response) {

            if(response.status === 200) {

                setupUserdata(response.responseText);
                dissmissAllModals();
                return true;

            }
            if(response.status === 405) {

                alert("Message under Construction ;)")
                dissmissAllModals();

            }
            if(response.status === 500) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }
            if(response.status === 401) {

                setupLoggedOutPage();
                unsetUserdata();

            }

        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();

        };

        xhr.send ();

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}


//////////////////////////////////////////////////////////////////////////////////////////////////////
//                                      Serverside Communication                                    //
//////////////////////////////////////////////////////////////////////////////////////////////////////
/**
 * 
 * @param {Event} event 
 */
function uploadImage(event) {

    try {

        //Get Formdata
        toggleLoadingScreen();
        event.preventDefault();
        var formElement = event.target;
        var formdata = new FormData (formElement)
        formdata.append("uploadtime", Date.now())
        var action = `/users/${applicationState.userID}/images`;
        var method = "POST";

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open(method, action, true);
        xhr.send(formdata);

        //Callbacks
        var success = function(response) {

            if(response.status===201) {

                location.reload(); 

            } 

            if(response.status===401) {

                alert("Something Made wrong: Message under Construction ;)");
                window.setTimeout(function(){dissmissAllModals();}, 1000);

            } 

            if(response.status === 409) {
                
                alert("Something Made wrong: Message under Construction ;)");
                window.setTimeout(function(){dissmissAllModals();}, 1000);

            }

            if(response.status === 500) {
                
                alert("Something Made wrong: Message under Construction ;)");
                window.setTimeout(function(){dissmissAllModals();}, 1000);

            }
            
        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            window.setTimeout(function(){dissmissAllModals();}, 1000);

        };   

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }




}document.getElementById("uploadImageForm").addEventListener('submit', function(ev){uploadImage(ev);})


/**
 * 
 * @param {Event} event 
 */
function commentImage(event) {

    try {

        //Get Formdata
        event.preventDefault();
        var formElement = event.target;
        var formdata    = new FormData (formElement)
        var imageID     = applicationState.currentCommentedImageID;
        var action      = `/images/${imageID}/comment`;
        var method      = "POST";

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open(method, action, true);

        //Callbacks
        var success = function(response) {

            if(response.status === 201) {

                dissmissAllModals();
                removeAllComments("image_" + applicationState.currentCommentedImageID);
                getComments(applicationState.currentCommentedImageID);

            }
            if(response.status >= 400) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }
            
        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();

        };   

        xhr.send(formdata);

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}document.getElementById("commentImageForm").addEventListener('submit', function(ev){commentImage(ev);})


/**
 * 
 * @param {string} imageID 
 */
function getComments(imageID) {

    try{

        //XHR
        var method = "GET";
        var action = `/images/${imageID}`;
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open (method, action, false);

        var success = function(response) {

            if(response.status === 200) {

                insertAllComments("image_" + imageID, response.responseText);

            }
            if(response.status === 404) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }
            if(response.status === 500) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }
        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();
        };

        xhr.send ();

    }catch{

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}


/**
 * 
 * @param {string} imageID 
 */
function getImageCard(imageID) {

    try{

        //XHR
        var method = "GET";
        var action = `/images/${imageID}`;
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open (method, action, false);

        var success = function(response) {

            if(response.status === 200) {

                var jsonObject = JSON.parse(response.responseText);
                addImageCard(response.responseText);
                getLike(jsonObject.ImageMetaData._id)
                showUnshowAllLikeButton();

            }
            if(response.status === 404) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }
            if(response.status === 500) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }
        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();
        };

        xhr.send ();

    }catch{

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}


/**
 * 
 * @param {number} recordTime 
 */
function getImageCards(recordTime) {

    try{

        var success = function(response) {

            if(response.status === 200) {

                var ids = JSON.parse(response.responseText).ImagesIDs;

                for (let index = 0; index < ids.length; index++) {

                    const id = ids[index];
                    getImageCard(id);

                }
            }
            if(response.status === 500) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }

        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();

        };

        //XHR
        var method = "GET";
        var action = `/images?lastRecordTime=${recordTime}`;
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open (method, action, false);
        xhr.send ();

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}


/**
 * 
 * @param {HTMLElement} cardElement 
 */
function likeImage(likeButton) {

    try {

        //Get Formdata
        var imageID = getImageCardID(likeButton);
        var action = `/images/${imageID}/like`;
        var method = "POST";

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open(method, action, true);
        xhr.send();

        //Callbacks
        var success = function(response) {

            if(response.status === 201) {

                var jsonObject = JSON.parse(response.responseText);

                changeLikeButtonLiked(jsonObject.ImageID)
                getLikes(jsonObject.ImageID);

            }
            if(response.status === 202) {

                var jsonObject = JSON.parse(response.responseText);

                changeLikeButtonNotLiked(jsonObject.ImageID)
                getLikes(jsonObject.ImageID);

            }

            if(response.status === 409) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }
            if(response.status === 500) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }
        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();

        };   

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}


/**
 * 
 * @param {number} imageID 
 */
function getLikes(imageID) {

    try {

        //Get Formdata
        var action = `/images/${imageID}`;
        var method = "GET";

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open(method, action, true);
        xhr.send();

        //Callbacks
        var success = function(response) {

            if(response.status === 200) {

                var jsonObject = JSON.parse(response.responseText);
                updateLikes(jsonObject.ImageMetaData._id, jsonObject.ImageMetaData.likes);

            }

            else {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }

        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();

        };   

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}


/**
 * 
 * @param {string} imageID 
 */
function getLike(imageID) {
    
    try {

        //Get Formdata
        var action = `/images/${imageID}/like`;
        var method = "GET";

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open(method, action, true);
        xhr.send();

        //Callbacks
        var success = function(response) {

            if(response.status === 200) {

                var jsonObject = JSON.parse(response.responseText);
                if(jsonObject.IsLiked) {
                    changeLikeButtonLiked(jsonObject.ImageID);
                }
                if(!jsonObject.IsLiked) {
                    changeLikeButtonNotLiked(jsonObject.ImageID);
                }

            }

            else {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }

        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();

        };   

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }
}

/**
 * 
 */
function getUserImages() {

    try{

        var success = function(response) {

            if(response.status === 200) {

                var images = JSON.parse(response.responseText).Images;

                for (let index = 0; index < images.length; index++) {

                    const image = images[index];
                    addUserImageCard(image);

                }
            }
            if(response.status === 500) {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }

        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();

        };

        //XHR
        var method = "GET";
        var action = `/users/${applicationState.userID}/images`;
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open (method, action, false);
        xhr.send ();


    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}

/**
 * 
 * @param {string} imageID 
 */
function deleteImage(imageID) {
    
    try {

        //Get Formdata
        var action = `/images/${imageID}`;
        var method = "DELETE";

        //XHR
        var xhr = new XMLHttpRequest();
        xhr.onload  = function(){ success(xhr); } // success case
        xhr.onerror = function(){   error(xhr); } // failure case
        xhr.open(method, action, true);
        xhr.send();

        //Callbacks
        var success = function(response) {

            if(response.status === 202) {

                removeAllUserImageCards();
                getUserImages();

            }

            else {

                alert("Something Made wrong: Message under Construction ;)")
                dissmissAllModals();

            }

        };
        var error   = function(response) {

            alert("Something Made wrong: Message under Construction ;)")
            dissmissAllModals();

        };   

    } catch {

        alert("Something Made wrong: Message under Construction ;)")
        dissmissAllModals();

    }

}


//////////////////////////////////////////////////////////////////////////////////////////////////////
//                                            UI Changer                                            //
//////////////////////////////////////////////////////////////////////////////////////////////////////

//Pages
/**
 * 
 */
function openMyImagesPage(){

    var publicPage = document.getElementById('publicImages');
    var usersPage  = document.getElementById( 'userImages' );

    publicPage.classList.add("hide");
    usersPage.classList.remove("hide");

}


/**
 * 
 */
function openPublicPage(){

    var publicPage = document.getElementById('publicImages');
    var usersPage  = document.getElementById( 'userImages' );

    publicPage.classList.remove("hide");
    usersPage.classList.add("hide");

}


// Like Changer
/**
 * 
 * @param {HTMLButtonElement} buttonElement 
 */
function changeLikeIcon(buttonElement) {

    buttonElement.children[0].classList = "far fa-heart";

}


/**
 * 
 * @param {HTMLButtonElement} buttonElement 
 */
function changeLikeIconHover(buttonElement) {

    buttonElement.children[0].classList = "fa fa-heart";

}


/**
 * 
 * @param {string} buttonElement 
 */
function changeLikeButtonLiked(imageID) {
    var imageCard = document.getElementById("image_" + imageID);
    var likeButton = imageCard.getElementsByClassName("likeButton")[0];
    likeButton.classList.add("likedButton")
}


/**
 * 
 * @param {string} buttonElement 
 */
function changeLikeButtonNotLiked(imageID) {
    var imageCard = document.getElementById("image_" + imageID);
    var likeButton = imageCard.getElementsByClassName("likeButton")[0];
    likeButton.classList.remove("likedButton")
}


/**
 * 
 */
function showUnshowAllLikeButton(){
    allLikeButtonContainer = document.getElementsByClassName("likeButtonContainer");

    for (let index = 0; index < allLikeButtonContainer.length; index++) {
        const container = allLikeButtonContainer[index];
        if(container.classList.contains("hide")&&applicationState.loggedIn){
            container.classList.remove("hide");
        }
    }

}


/**
 * 
 * @param {string} imageID 
 * @param {number} likes 
 */
function updateLikes(imageID, likes) {

    var imageCard = document.getElementById("image_" + imageID);
    var likeCounter = imageCard.getElementsByClassName("likeCounter")[0];
    likeCounter.innerHTML = likes || 0;

}


//Comments
/**
 * 
 * @param {string} imageID 
 */
function showAllComments(imageID) {

    //TODO

}


/**
 * 
 * @param {string} imageID 
 * @param {JSON} jsonObject 
 */
function insertAllComments(imageID, jsonObject){

    var responseJSON = JSON.parse(jsonObject);

    var imageCard = document.getElementById(imageID);
    var commentArea = imageCard.getElementsByClassName("commentArea")[0];
    var commentsCounter = imageCard.getElementsByClassName('commentsCounter')[0];
    commentsCounter.innerHTML = responseJSON.Comments.length;


    if(responseJSON.Comments !== null) {

        for (let index = 0; index < responseJSON.Comments.length; index++) {
        
            const comment = responseJSON.Comments[index];
    
            //Get and Clone Template
            var commentTemp  = commentTemplate.content.cloneNode(true).children[0];
            var commentTempclone = commentTemp.cloneNode(true);
    
            var commentUsername = commentTempclone.getElementsByClassName('commentUsername')[0];
            var commentText     = commentTempclone.getElementsByClassName('commentText')[0];
            
            commentUsername.innerHTML = comment.owner;
            commentUsername.title = comment.owner;
            commentText.innerHTML = comment.comment;
    
            commentArea.appendChild(commentTempclone)
    
        }

    }
}


/**
 * 
 * @param {string} imageID 
 */
function removeAllComments(imageID) {

    var imageCard = document.getElementById(imageID);
    var comments  = imageCard.getElementsByClassName("commentArea")[0];

    comments.innerHTML = "";

}


/**
 * 
 * @param {HTMLElement} htmlElement 
 */
function setCurrentCommentedImageID(htmlElement) {
    var startChild = htmlElement;
    if(startChild.id.includes("image_"))
        applicationState.currentCommentedImageID = imageID = startChild.id.substring(6, rawID.length);//The first 6 chars are "image_" the rest is the id
    for (let p = startChild.parentElement; p.tagName != "BODY"; p = p.parentElement) {
        if(p.id.includes("image_"))
            applicationState.currentCommentedImageID = imageID = p.id.substring(6, p.length);//The first 6 chars are "image_" the rest is the id
    }
}


//Images
/**
 * 
 * @param {HTMLElement} htmlElement 
 */
function setCurrentImageToDeleteID(htmlElement) {
    var startChild = htmlElement;
    if(startChild.id.includes("userImage_"))
        return applicationState.currentCommentedImageID = imageID = startChild.id.substring(10, rawID.length);//The first 6 chars are "image_" the rest is the id
    for (let p = startChild.parentElement; p.tagName != "BODY"; p = p.parentElement) {
        if(p.id.includes("userImage_"))
        return applicationState.currentCommentedImageID = imageID = p.id.substring(10, p.length);//The first 6 chars are "image_" the rest is the id
    }
}


/**
 * 
 */
function loadUserImages() {
    //TODO
}


/**
 * 
 * @param {string} imageSrc 
 * @param {number} likeCount 
 * @param {number} commentCount 
 */
function addUserImageCard(jsonObject){
    
    //Get Json
    var even = (applicationState.userimagesCount%2 === 0);

    if(!even) { //If a the new amount of new Images is uneven

        //Get and Clone Template
        var template  = userImageCardTemplate.content.cloneNode(true).children[0];
        var tempclone = template.cloneNode(true);
        var dummyCard = userImages.getElementsByClassName("userImageDummy")[0];

        //Get Elements in Template
        var imageElement   = tempclone.getElementsByClassName(     'userImageCardImage'    )[0];
        var imageUpload    = tempclone.getElementsByClassName(   'userImageCardDateTime'   )[0];
        var likeCounter    = tempclone.getElementsByClassName(  'userImageCardLikeCounter' )[0];
        var commentCounter = tempclone.getElementsByClassName('userImageCardCommentCounter')[0];

        //Fill Object
        imageElement.src         = `../images/${jsonObject._id}/raw`;
        likeCounter.innerHTML    = jsonObject.likes || 0;
        commentCounter.innerHTML = jsonObject.comments || 0;
        imageUpload.innerHTML    = new Date(jsonObject.uploadtime).toTimeString();

        //Append written Object
        dummyCard.parentElement.replaceChild(tempclone,dummyCard);
        applicationState.userimagesCount++;

    } else {    //If a the new amount of new Images is even
        
        //Get and Clone Template
        var template  = userImageCardRowTemplate.content.cloneNode(true).children[0];
        var tempclone = template.cloneNode(true);

        //Get Elements in Template
        var imageCard      = tempclone.getElementsByClassName(       'userImageCard'       )[0];
        var imageElement   = tempclone.getElementsByClassName(     'userImageCardImage'    )[0];
        var imageUpload    = tempclone.getElementsByClassName(   'userImageCardDateTime'   )[0];
        var likeCounter    = tempclone.getElementsByClassName(  'userImageCardLikeCounter' )[0];
        var commentCounter = tempclone.getElementsByClassName('userImageCardCommentCounter')[0];

        //Fill Object
        imageElement.src         = `../images/${jsonObject._id}/raw`;
        imageCard.id             = `userImage_${jsonObject._id}`;
        likeCounter.innerHTML    = jsonObject.likes || 0;
        commentCounter.innerHTML = jsonObject.comments || 0;
        imageUpload.innerHTML    = new Date(jsonObject.uploadtime).toTimeString();

        //Appen written Object
        userImages.appendChild(tempclone);
        applicationState.userimagesCount++;

    }

}


/**
 * 
 * @param {string} imageID 
 */
function removeAllUserImageCards() {

    var userImageCards = userImages.getElementsByClassName("userImageCardRow");
    var length = userImageCards.length;

    for (let index = 0; index < length; index++) {
        const userImageCard = userImageCards[0];

        userImageCard.remove();
        
    }

    applicationState.userimagesCount = 0;

}


/**
 * 
 * @param {string} jsonObject 
 */
function addImageCard(jsonObject) {

    //Get Json
    var responseJSON = JSON.parse(jsonObject)

    //Get and Clone Template
    var template  = imageCardTemplate.content.cloneNode(true).children[0];
    var tempclone = template.cloneNode(true);

    //Get Elements in Template
    tempclone.id            = `image_${responseJSON.ImageMetaData._id}`;
    var imageElement        = tempclone.getElementsByClassName(    'imageElement'   )[0];
    var likeCounter         = tempclone.getElementsByClassName(     'likeCounter'   )[0];
    var imageDescription    = tempclone.getElementsByClassName(  'imageDescription' )[0];
    var userName            = tempclone.getElementsByClassName(      'userName'     )[0];
    var datetime            = tempclone.getElementsByClassName(      'datetime'     )[0];
    var likeButton          = tempclone.getElementsByClassName(     'likeButton'    )[0];
    var likeButtonContainer = tempclone.getElementsByClassName('likeButtonContainer')[0];
    var commentArea         = tempclone.getElementsByClassName(    'commentArea'    )[0];
    var commentsCounter     = tempclone.getElementsByClassName(  'commentsCounter'  )[0];
    
    //Fill Object TEST
    imageElement.src           = `/images/${responseJSON.ImageMetaData._id}/raw`;
    likeCounter.innerHTML      = responseJSON.ImageMetaData.likes || 0;
    imageDescription.innerHTML = responseJSON.ImageMetaData.description;
    userName.innerHTML         = responseJSON.ImageMetaData.owner;
    datetime.innerHTML         = new Date(responseJSON.ImageMetaData.uploadtime).toTimeString()
    commentsCounter.innerHTML  = 0;

    if(responseJSON.Comments !== null) {

        commentsCounter.innerHTML  = responseJSON.Comments.length | 0;


        for (let index = 0; index < responseJSON.Comments.length; index++) {
        
            const comment = responseJSON.Comments[index];
    
            //Get and Clone Template
            var commentTemp  = commentTemplate.content.cloneNode(true).children[0];
            var commentTempclone = commentTemp.cloneNode(true);
    
            var commentUsername = commentTempclone.getElementsByClassName('commentUsername')[0];
            var commentText     = commentTempclone.getElementsByClassName('commentText')[0];
            
            commentUsername.innerHTML = comment.owner;
            commentUsername.title = comment.owner;
            commentText.innerHTML = comment.comment;
    
            commentArea.appendChild(commentTempclone)
    
        }

    }


    //Appen written Object
    publicImages.appendChild(tempclone);

    applicationState.lastimageDateTime = responseJSON.ImageMetaData.uploadtime
}


/**
 * 
 * @param {HTMLElement} htmlElement 
 */
function getImageCardID(htmlElement) {
    var startChild = htmlElement;
    if(startChild.id.includes("image_"))
        return imageID = startChild.id.substring(5, rawID.length);//The first 6 chars are "image_" the rest is the id
    for (let p = startChild.parentElement; p.tagName != "BODY"; p = p.parentElement) {
        if(p.id.includes("image_"))
            return imageID = p.id.substring(6, p.length);//The first 6 chars are "image_" the rest is the id
    }
}


//Sessionmanagement
/**
 * 
 */
function setupLoggedPage() {

    //Change Url
    window.history.pushState(null, `${applicationState.userName} - Flashlight`, "/users/" + applicationState.userID);

    //Get Elements
    likeButtons    = publicImages.getElementsByClassName("likeButtonContainer"); 
    commentAreas   = publicImages.getElementsByClassName("commenArea");
    commentButtons = publicImages.getElementsByClassName("commentButton");

    //ChangeUI Template
    headerLoggedIn.classList.remove("hide");
    headerLoggedOut.classList.add("hide");

    //Activate all LikeButton from all Loadedimages
    for (let index = 0; index < likeButtons.length; index++) {     
        const likeButton = likeButtons[index];
        likeButton.classList.remove("hide");
    }

    //Activate all Comments from all Loadedimages
    for (let index = 0; index < commentAreas.length; index++) {
        const commentArea = commentAreas[index];
        commentArea.classList.remove("hide");
    }

    //Activate all Commentbuttons from all Loadedimages
    for (let index = 0; index < commentButtons.length; index++) {
        const commentButton = commentButtons[index];
        commentButton.classList.remove("hide");
    }
}


/**
 * 
 */
function setupLoggedOutPage() {

    //Change Url
    window.history.pushState(null, "Flashlight - By Asef Alper Tunga Dündar", "/");

    //Get all Elements
    likeButtons    = publicImages.getElementsByClassName("likeButtonContainer"); 
    commentAreas   = publicImages.getElementsByClassName("commenArea");
    commentButtons = publicImages.getElementsByClassName("commentButton");

    //ChangeUI Template
    headerLoggedIn.classList.add("hide");
    headerLoggedOut.classList.remove("hide");

    //Deavtivate all LikeButton from all Loadedimages
    for (let index = 0; index < likeButtons.length; index++) {     
        const likeButton = likeButtons[index];
        likeButton.classList.add("hide");
    }

    //Deavtivate all Commentareas from all Loadedimages
    for (let index = 0; index < commentAreas.length; index++) {     
        const commentArea = commentAreas[index];
        commentArea.classList.add("hide");
    }

    //Deactivate all Commentbuttons from all Loadedimages
    for (let index = 0; index < commentButtons.length; index++) {
        const commentButton = commentButtons[index];
        commentButton.classList.add("hide");
    }

}


/**
 * 
 * @param {string} jsonText 
 */
function setupUserdata(jsonText) {

    var jsonParsed = JSON.parse(jsonText);
    applicationState.userID   = jsonParsed.HashedUsername;
    applicationState.username = jsonParsed.Username;
    applicationState.loggedIn = true;

    loggedInName.innerHTML = applicationState.username;
    loggedInNameMyImages.innerHTML = applicationState.userName;

    getUserImages();
    
}


/**
 * 
 */
function initRecordTimes() {

    applicationState.lastimageDateTime  = Date.now();
    applicationState.firstImageDataTime = Date.now();

}


/**
 * 
 */
function unsetUserdata() {

    applicationState.userID          = "";
    applicationState.userName        = "";
    applicationState.userimagesCount = 0;

    document.getElementById("loggedInName").innerHTML         = "";
    document.getElementById("loggedInNameMyImages").innerHTML = "";

}

//Modals
/**
 * 
 */
function dissmissAllModals() {

    var modals = document.getElementsByClassName("modal");
    for (let index = 0; index < modals.length; index++) {
        const modalElement = modals[index];
        $('#'+modalElement.id).modal('hide');
    }

}


/**
 * 
 */
function toggleLoadingScreen(){

    $('#loadingscreen').modal('toggle');

}


/**
 * 
 */
function unshowLoadingScreen() {

    $('#loadingscreen').modal({show:false});

}


//////////////////////////////////////////////////////////////////////////////////////////////////////
//                                               Setup                                              //
//////////////////////////////////////////////////////////////////////////////////////////////////////
//Set Loadingscreenmodal and make escaping impossible
$('#loadingscreen').modal({backdrop: 'static', keyboard: false, show:false});

//ApplicationState Context
var applicationState = {
    userName:                "",
    userID:                  "",
    loggedIn:                false,
    lastimageDateTime:       "",
    firstImageDataTime:      "",
    userimagesCount:         0,
    currentCommentedImageID: "",
}

//Templates and Static HTMLElement
var publicImages              = document.getElementById(       "publicImages"      );
var userImages                = document.getElementById(        "userImages"       );
var headerLoggedIn            = document.getElementById(        'loggedBar'        );
var headerLoggedOut           = document.getElementById(       'notLoggedBar'      );
var loggedInName              = document.getElementById(       "loggedInName"      );
var loggedInNameMyImages      = document.getElementById(   "loggedInNameMyImages"  );
var imageCardTemplate         = document.getElementById(     "imageCardTemplate"   );
var commentTemplate           = document.getElementById(      "commentTemplate"    );
var userImageCardRowTemplate  = document.getElementById("usersImageCardRowTemplate");
var userImageCardTemplate     = document.getElementById(  "userImageCardTemplate"  );

window.onload   = function() {

    
    initRecordTimes();

    getUserdata()
    if(applicationState.loggedIn){
        setupLoggedPage();
    } 

    getImageCards(applicationState.lastimageDateTime);

};

window.onscroll = function() {
    if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight && publicImages.children.length < 100) {
        getImageCards(applicationState.lastimageDateTime);
    }
};