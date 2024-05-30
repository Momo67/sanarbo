// ComponentToTest.spec.js
import { describe, it, expect} from 'vitest';
import { mount, config } from '@vue/test-utils';
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import ComponentToTest from '../ComponentToTest.vue';

const vuetify = createVuetify({
  components,
  directives,
})
config.global.plugins = [vuetify]

global.ResizeObserver = require('resize-observer-polyfill')

describe('ComponentToTest', () => {
  it('emits button-clicked event when button is clicked', async () => {
    const wrapper = mount(ComponentToTest);
    const button = wrapper.find('button');

    await button.trigger('click');

    expect(wrapper.emitted('button-clicked')).toBeTruthy();
  });

  it('renders the button with text "Click me"', () => {
    const wrapper = mount(ComponentToTest);
    const button = wrapper.find('button');

    expect(button.text()).toBe('Click me');
  });

  it('renders VBtn from Vuetify', () => {
    const wrapper = mount(ComponentToTest);
    const vButton = wrapper.findComponent({ name: 'VBtn' });

    expect(vButton.exists()).toBe(true);
  });
});
