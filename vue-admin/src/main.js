// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store'
import instance from './axios.js'
import VeeValidate, { Validator } from 'vee-validate'
import ES from 'vee-validate/dist/locale/es'

import {
  Vuetify,
  VApp,
  VAlert,
  VAutocomplete,
  VAvatar,
  VBadge,
  VBottomNav,
  VBottomSheet,
  VBreadcrumbs,
  VBtn,
  VBtnToggle,
  VCard,
  VCarousel,
  VCheckbox,
  VChip,
  VCombobox,
  VCounter,
  VDataIterator,
  VDataTable,
  VDatePicker,
  VDialog,
  VDivider,
  VExpansionPanel,
  VFooter,
  VForm,
  VGrid,
  VIcon,
  VInput,
  VJumbotron,
  VLabel,
  VList,
  VMenu,
  VMessages,
  VNavigationDrawer,
  VOverflowBtn,
  VPagination,
  VParallax,
  VProgressCircular,
  VProgressLinear,
  VRadioGroup,
  VRangeSlider,
  VSelect,
  VSlider,
  VSnackbar,
  VSpeedDial,
  VStepper,
  VSubheader,
  VSwitch,
  VSystemBar,
  VTabs,
  VTextField,
  VTextarea,
  VTimePicker,
  VToolbar,
  VTooltip,
  VImg,
  transitions
} from 'vuetify'
import resize from 'vuetify/es5/directives/resize'
import '../node_modules/vuetify/src/stylus/app.styl'

Validator.localize('es', ES)
Vue.use(VeeValidate)

Vue.use(Vuetify, {
  components: {
    VApp,
    VAlert,
    VAutocomplete,
    VAvatar,
    VBadge,
    VBottomNav,
    VBottomSheet,
    VBreadcrumbs,
    VBtn,
    VBtnToggle,
    VCard,
    VCarousel,
    VCheckbox,
    VChip,
    VCombobox,
    VCounter,
    VDataIterator,
    VDataTable,
    VDatePicker,
    VDialog,
    VDivider,
    VExpansionPanel,
    VFooter,
    VForm,
    VGrid,
    VIcon,
    VInput,
    VJumbotron,
    VLabel,
    VList,
    VMenu,
    VMessages,
    VNavigationDrawer,
    VOverflowBtn,
    VPagination,
    VParallax,
    VProgressCircular,
    VProgressLinear,
    VRadioGroup,
    VRangeSlider,
    VSelect,
    VSlider,
    VSnackbar,
    VSpeedDial,
    VStepper,
    VSubheader,
    VSwitch,
    VSystemBar,
    VTabs,
    VTextField,
    VTextarea,
    VTimePicker,
    VToolbar,
    VTooltip,
    VImg,
    transitions
  },
  directives: {
    resize: resize
  },
  theme: {
    primary: '#3D5AFE',
    secondary: '#FAFAFA',
    accent: '#F4511E',
    error: '#D50000',
    warning: '#ffeb3b',
    info: '#2196f3',
    success: '#4caf50',
    dark: '#3d3c3c'
  }
})

// Add a request interceptor
instance.interceptors.request.use(function (config) {
  if (location.pathname !== '/') {
    const token = localStorage.getItem('bazam-token')
    if (token !== null && token !== undefined && token !== '') {
      config.headers['Authorization'] = token
    } else {
      localStorage.removeItem('bazam-token')
      location.href = '/'
    }
  }
  return config
})

instance.interceptors.response.use(function (response) {
  return response
}, function (error) {
  if (error.response) {
    // console.log(error.response)
    // console.log(error.response.status);
    // console.log(error.response.headers);
    switch (error.response.status) {
      case 401:
        localStorage.clear()
        location.reload()
        break
      case 409:
        alert('El elemento esta uso, por favor liberalo antes de eliminarlo')
        break
      default:
        console.log('Server response: ', error.response.data.pretty_message)
        break
    }
  } else if (error.request) {
    console.log('Request error: ', error.request)
  } else {
    console.log('Error', error.message)
  }
})

Vue.config.productionTip = false
Vue.prototype.$axios = instance

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
