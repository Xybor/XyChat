(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-2d0d6d35"],{"73cf":function(e,t,o){"use strict";o.r(t);var s=o("7a23"),r={class:"login container mt-5"},c=Object(s["g"])('<h3 class="text-center">Register</h3><div class="row justify-content-md-center"><div class="col-3"><div class="form-group"><label for="firstname">First name</label><input type="text" class="form-control" id="firstname"></div></div><div class="col-3"><div class="form-group"><label for="lastname">Last name</label><input type="text" class="form-control" id="lastname"></div></div></div>',2),a={class:"row justify-content-md-center"},n={class:"form-group col-6"},l=Object(s["i"])("label",{for:"username"},"Username",-1),i=Object(s["g"])('<div class="row justify-content-md-center mt-2"><div class="form-group mt-2 col-6"><label for="exampleFormControlInput1">Gender</label><select class="form-select" aria-label="Default select example"><option value="1" selected>Male</option><option value="2">Female</option><option value="3">Other</option></select></div></div>',1),d={class:"row justify-content-md-center mt-2"},u={class:"form-group col-6"},p=Object(s["i"])("label",{for:"password"},"Password",-1),m={class:"row justify-content-md-center mt-2"},b={class:"form-group col-6"},f=Object(s["i"])("label",{for:"password"},"Verify Password",-1),j={class:"row justify-content-md-center mt-2"},v={class:"form-group col-6"},O=Object(s["h"])(" Switch to "),w=Object(s["h"])("Login"),y={class:"text-center mt-3"};function h(e,t,o,h,g,x){var U=Object(s["B"])("router-link");return Object(s["t"])(),Object(s["e"])("div",r,[c,Object(s["i"])("div",a,[Object(s["i"])("div",n,[l,Object(s["L"])(Object(s["i"])("input",{type:"text",class:"form-control",id:"username",placeholder:"Username","onUpdate:modelValue":t[1]||(t[1]=function(e){return h.username=e})},null,512),[[s["I"],h.username]])])]),i,Object(s["i"])("div",d,[Object(s["i"])("div",u,[p,Object(s["L"])(Object(s["i"])("input",{type:"password",class:"form-control",id:"password",placeholder:"******","onUpdate:modelValue":t[2]||(t[2]=function(e){return h.password=e})},null,512),[[s["I"],h.password]])])]),Object(s["i"])("div",m,[Object(s["i"])("div",b,[f,Object(s["L"])(Object(s["i"])("input",{type:"password",class:"form-control",id:"password",placeholder:"******","onUpdate:modelValue":t[3]||(t[3]=function(e){return h.password=e})},null,512),[[s["I"],h.password]])])]),Object(s["i"])("div",j,[Object(s["i"])("div",v,[O,Object(s["i"])(U,{to:"/login"},{default:Object(s["K"])((function(){return[w]})),_:1})])]),Object(s["i"])("div",y,[Object(s["i"])("button",{class:"btn btn-outline-success text-center",type:"submit",onClick:t[4]||(t[4]=function(){return h.handleSubmit&&h.handleSubmit.apply(h,arguments)})}," Create new account ")])])}var g=o("5502"),x={setup:function(){var e=Object(s["y"])(""),t=Object(s["y"])(""),o=Object(g["b"])(),r=function(){if(""==e.value||""==t.value)return o.dispatch("alert/error","Username or password can not empty"),!1;o.dispatch("account/register",{username:e.value,password:t.value})};return{username:e,password:t,handleSubmit:r}}};x.render=h;t["default"]=x}}]);
//# sourceMappingURL=chunk-2d0d6d35.43024665.js.map