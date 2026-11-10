-- object_validator is checked by the trees service (ValidationList, SaveValidation)
-- but was missing from 000005. ON CONFLICT keeps databases where it was created by hand untouched.
INSERT INTO public.go_group (name, creator, comment)
VALUES ('object_validator', 1, 'validators of object microservice')
ON CONFLICT (name) DO NOTHING;
