import HomeComp from './components'
import Vue from "vue"
var vm = new Vue({
    el: "#app",
    components : {
       home : HomeComp
    }
})