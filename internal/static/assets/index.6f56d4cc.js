import{z as t,P as I,ck as b,r,aG as y,ac as n,cl as c,cm as f,u as m,o as h,w,ch as P,a5 as S,bA as F}from"./index.2b950906.js";import{a as _}from"./size.4b7247cf.js";const z=u=>{const s=I();return t(()=>{var o,e;return(e=((o=s.proxy)==null?void 0:o.$props)[u])!=null?e:void 0})},U=b({type:String,values:_,required:!1}),g=(u,s={})=>{const o=r(void 0),e=s.prop?o:z("size"),l=s.global?o:y("size"),a=s.form?{size:void 0}:n(c,void 0),v=s.formItem?{size:void 0}:n(f,void 0);return t(()=>e.value||m(u)||(v==null?void 0:v.size)||(a==null?void 0:a.size)||l.value||"")},j=u=>{const s=z("disabled"),o=n(c,void 0);return t(()=>s.value||m(u)||(o==null?void 0:o.disabled)||!1)},k=()=>{const u=n(c,void 0),s=n(f,void 0);return{form:u,formItem:s}},q=(u,{formItemContext:s,disableIdGeneration:o,disableIdManagement:e})=>{o||(o=r(!1)),e||(e=r(!1));const l=r();let a;const v=t(()=>{var i;return!!(!u.label&&s&&s.inputIds&&((i=s.inputIds)==null?void 0:i.length)<=1)});return h(()=>{a=w([S(u,"id"),o],([i,p])=>{const d=i!=null?i:p?void 0:P().value;d!==l.value&&(s!=null&&s.removeInputId&&(l.value&&s.removeInputId(l.value),!(e!=null&&e.value)&&!p&&d&&s.addInputId(d)),l.value=d)},{immediate:!0})}),F(()=>{a&&a(),s!=null&&s.removeInputId&&l.value&&s.removeInputId(l.value)}),{isLabeledByFormItem:v,inputId:l}};export{k as a,q as b,g as c,j as d,U as u};