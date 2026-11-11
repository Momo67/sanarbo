import { mount, flushPromises, config } from '@vue/test-utils';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ref } from 'vue';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';
import TreeForm from '../TreeForm.vue';
import { useFetch } from '../../composables/FetchData.js';

const vuetify = createVuetify({
  components,
  directives,
})
config.global.plugins = [vuetify]

global.ResizeObserver = require('resize-observer-polyfill')

const storedTree = {
  id: 42,
  external_id: 1042,
  name: 'Platane 42',
  description: 'Avenue de la Gare',
  is_active: true,
  is_validated: true,
  create_time: '2024-01-01T00:00:00Z',
  creator: 1007,
  comment: 'arbre remarquable',
  inactivation_reason: 'ancienne raison',
  geom: 'POINT(2538000 1152000)',
  tree_attributes: { idvalidation: 1, circonference: 120 },
};

vi.mock('../../composables/FetchData.js', () => ({
  useFetch: vi.fn(),
}));

// fetchResult mimics what useFetch returns
const fetchResult = (data, status = 200, errorMessage = '') => ({
  data: ref(data),
  status: ref(status),
  hasError: ref(status < 200 || status > 299),
  errorMessage: ref(errorMessage),
  isLoading: ref(false),
});

// mockBackend answers the GET with storedTree and the PUT with putResult
const mockBackend = (putResult) => {
  useFetch.mockImplementation(async (url, options) => {
    if (options.method === 'PUT') {
      return putResult ?? fetchResult(JSON.parse(options.body));
    }
    return fetchResult(structuredClone(storedTree));
  });
};

const mountForm = async (mountOptions = {}) => {
  const dictionaries = new Proxy({ validation: { data: [{ id: 1, value: 'à valider' }] } }, {
    get: (target, name) => target[name] ?? { data: [] },
  });
  const wrapper = mount(TreeForm, { props: { showForm: true, treeId: 42, dictionaries }, ...mountOptions });
  await flushPromises();
  return wrapper;
};

const putCalls = () => useFetch.mock.calls.filter(([, options]) => options.method === 'PUT');

describe('TreeForm', () => {
  beforeEach(() => {
    useFetch.mockReset();
  });

  it('sends back the fields it does not edit when saving', async () => {
    mockBackend();
    const wrapper = await mountForm();

    await wrapper.vm.submitForm();

    expect(putCalls()).toHaveLength(1);
    const body = JSON.parse(putCalls()[0][1].body);
    expect(body.comment).toBe('arbre remarquable');
    expect(body.inactivation_reason).toBe('ancienne raison');
    expect(body.is_validated).toBe(false);
  });

  it('saves once and closes when the save button is clicked', async () => {
    mockBackend();
    // jsdom only submits a form that is attached to the document
    const wrapper = await mountForm({ attachTo: document.body });

    await wrapper.find('button[type="submit"]').trigger('click');
    await flushPromises();

    expect(putCalls()).toHaveLength(1);
    expect(wrapper.emitted('formSubmitted')).toHaveLength(1);
    wrapper.unmount();
  });

  it.each([
    [403, 'current user has no Goeland user id (external_id), cannot create or modify a tree', 'identifiant Goéland'],
    [403, 'current user has no editor privilege on this object', 'pas le droit de modifier cet arbre'],
    [401, 'invalid or expired jwt', 'session a expiré'],
    [500, 'Update got problem updating tree', 'erreur 500'],
  ])('shows the error and stays open when the backend answers %i', async (status, serverMessage, shown) => {
    mockBackend(fetchResult(null, status, serverMessage));
    const wrapper = await mountForm();

    await wrapper.vm.submitForm();
    await flushPromises();

    expect(wrapper.emitted('formSubmitted'), 'a failed save must not be reported as done').toBeUndefined();
    expect(wrapper.text()).toContain(shown);
  });
});
