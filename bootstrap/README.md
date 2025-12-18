# Bootstrap

The pre-requisites that has to be applied manually WITHOUT storing the state.
Instead all resources has to have corresponding import statement.

1. Add terraform code.
2. Apply terraform changes.
3. Add new/modified resources to imports.tf.
4. Remove local terraform state.
5. Run terraform apply to verify.
