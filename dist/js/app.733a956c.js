(function(){"use strict";var t={2043:function(t,n,e){var o=e(6369),i=function(){var t=this,n=t._self._c;return n("div",{attrs:{id:"app"}},[n("div",{staticClass:"my-input"},[n("el-form",{attrs:{"label-position":"top","label-width":"80px",model:t.config}},[n("el-form-item",{attrs:{label:"Cookie"}},[n("el-input",{model:{value:t.config.cookie,callback:function(n){t.$set(t.config,"cookie",n)},expression:"config.cookie"}})],1),n("el-form-item",{attrs:{label:"CSRF Token"}},[n("el-input",{model:{value:t.config.token,callback:function(n){t.$set(t.config,"token",n)},expression:"config.token"}})],1),n("el-form-item",{attrs:{label:"需要多少天内的日报"}},[n("el-input",{model:{value:t.config.length,callback:function(n){t.$set(t.config,"length",t._n(n))},expression:"config.length"}})],1)],1),n("el-button",{attrs:{type:"primary",loading:t.isLoading},on:{click:t.submit}},[t._v("搜索")])],1),n("div",{staticClass:"my-table"},[n("el-table",{attrs:{data:t.results,height:"450"}},[n("el-table-column",{attrs:{prop:"sendAt",label:"发送时间",width:"150",sortable:""}}),n("el-table-column",{attrs:{prop:"content",label:"内容"}}),n("el-table-column",{attrs:{fixed:"right",label:"操作",width:"100"},scopedSlots:t._u([{key:"default",fn:function(e){return[n("el-button",{attrs:{type:"primary",size:"small"},on:{click:function(n){return t.copyToClipboard(e.row.content)}}},[t._v("复制")])]}}])})],1)],1)])},r=[],l=e(2140),a=e.n(l),c=e(1168),s=e.n(c),u=e(8319),f=e.n(u),p=e(1540),d=e.n(p),b=e(5981),h=e.n(b),g=e(3480),m=e.n(g),v=e(2086),k=e.n(v),y=e(2482),w=e(6265),O=e.n(w);class x{}(0,y.Z)(x,"fetch",(t=>O().post("/mails/search",t)));var T={name:"App",components:{ElTable:k(),ElTableColumn:m(),ElInput:h(),ElButton:d(),ElForm:f(),ElFormItem:s()},data(){return{results:[],config:{length:0,cookie:"",token:""},isLoading:!1}},methods:{async submit(){this.isLoading=!0,this.results=[];const t={length:this.config.length,cookie:this.config.cookie,token:this.config.token};try{const n=await x.fetch(t),e=n.data.items;window.open(n.data.fileURL),e.forEach((t=>{this.results.push(t)}))}catch(n){a().error(n.response.data)}this.isLoading=!1},copyToClipboard(t){navigator.clipboard.writeText(t),a().success("copied!")}}},_=T,C=e(1001),j=(0,C.Z)(_,i,r,!1,null,null,null),E=j.exports;o["default"].config.productionTip=!1,new o["default"]({render:t=>t(E)}).$mount("#app")}},n={};function e(o){var i=n[o];if(void 0!==i)return i.exports;var r=n[o]={exports:{}};return t[o](r,r.exports,e),r.exports}e.m=t,function(){var t=[];e.O=function(n,o,i,r){if(!o){var l=1/0;for(u=0;u<t.length;u++){o=t[u][0],i=t[u][1],r=t[u][2];for(var a=!0,c=0;c<o.length;c++)(!1&r||l>=r)&&Object.keys(e.O).every((function(t){return e.O[t](o[c])}))?o.splice(c--,1):(a=!1,r<l&&(l=r));if(a){t.splice(u--,1);var s=i();void 0!==s&&(n=s)}}return n}r=r||0;for(var u=t.length;u>0&&t[u-1][2]>r;u--)t[u]=t[u-1];t[u]=[o,i,r]}}(),function(){e.n=function(t){var n=t&&t.__esModule?function(){return t["default"]}:function(){return t};return e.d(n,{a:n}),n}}(),function(){e.d=function(t,n){for(var o in n)e.o(n,o)&&!e.o(t,o)&&Object.defineProperty(t,o,{enumerable:!0,get:n[o]})}}(),function(){e.g=function(){if("object"===typeof globalThis)return globalThis;try{return this||new Function("return this")()}catch(t){if("object"===typeof window)return window}}()}(),function(){e.o=function(t,n){return Object.prototype.hasOwnProperty.call(t,n)}}(),function(){e.r=function(t){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})}}(),function(){var t={143:0};e.O.j=function(n){return 0===t[n]};var n=function(n,o){var i,r,l=o[0],a=o[1],c=o[2],s=0;if(l.some((function(n){return 0!==t[n]}))){for(i in a)e.o(a,i)&&(e.m[i]=a[i]);if(c)var u=c(e)}for(n&&n(o);s<l.length;s++)r=l[s],e.o(t,r)&&t[r]&&t[r][0](),t[r]=0;return e.O(u)},o=self["webpackChunkmailparser"]=self["webpackChunkmailparser"]||[];o.forEach(n.bind(null,0)),o.push=n.bind(null,o.push.bind(o))}();var o=e.O(void 0,[998],(function(){return e(2043)}));o=e.O(o)})();
//# sourceMappingURL=app.733a956c.js.map