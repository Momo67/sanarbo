import {reactive, toRefs} from "vue";

export const useFetch = async (url, options) => {
    const state = reactive(
        {
            isLoading: true,
            hasError: false,
            errorMessage: '',
            // HTTP status of the response, 0 when the request did not reach the server
            status: 0,
            data: [],
        }
    )



    const fetchData = async () => {
        state.isLoading = true;

        try {
            const res = await fetch(url, options);
            state.status = res.status;

            if (!res.ok) {
                // the backend answers errors with a JSON string or a {"message"} object
                let message = '';
                try {
                    const body = await res.json();
                    message = typeof body === 'string' ? body : (body?.message ?? '');
                } catch {
                    // body is not JSON
                }
                throw new Error(message || res.statusText || `HTTP ${res.status}`);
            }

            state.data = await res.json();

        } catch (err) {
            const typedError = err;
            state.hasError = true;
            state.errorMessage = typedError.message;
        } finally {
            state.isLoading = false;
        }
    };

        await fetchData();

        return {...toRefs(state)};
    }
