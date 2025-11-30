# Implementation Plan

## 1. Store Supplier Product Binding

- [ ] 1.1 Enhance StoreSupplierModule with first-bind default logic
  - Update `BindProducts` to set `is_default=true` for first binding of a product name
  - Add logic to check if same-named product already exists for the store
  - _Requirements: 2.3_

- [ ] 1.2 Write property test for binding idempotency
  - **Property 1: Binding Idempotency**
  - **Validates: Requirements 1.2**

- [ ] 1.3 Write property test for default supplier exclusivity
  - **Property 3: Default Supplier Exclusivity**
  - **Validates: Requirements 2.2**

- [ ] 1.4 Add ValidateStoreProducts method to StoreSupplierModule
  - Implement method to check if product IDs are bound to a store
  - Return list of unbound product IDs for error reporting
  - _Requirements: 3.4_

- [ ] 1.5 Write property test for unbind removes record
  - **Property 2: Unbind Removes Record**
  - **Validates: Requirements 1.3**

## 2. Purchase Order Creation Enhancement

- [ ] 2.1 Update PurchaseOrderService.CreateOrder with validation
  - Add call to `ValidateStoreProducts` before creating order
  - Return detailed error with unbound product IDs if validation fails
  - _Requirements: 3.4_

- [ ] 2.2 Write property test for unbound product rejection
  - **Property 7: Unbound Product Rejection**
  - **Validates: Requirements 3.4**

- [ ] 2.3 Write property test for supplier auto-matching
  - **Property 5: Supplier Auto-Matching**
  - **Validates: Requirements 3.2**

- [ ] 2.4 Write property test for total amount calculation
  - **Property 6: Total Amount Calculation Invariant**
  - **Validates: Requirements 3.3**

- [ ] 2.5 Write property test for order number uniqueness
  - **Property 4: Order Number Uniqueness**
  - **Validates: Requirements 3.1**

- [ ] 2.6 Write property test for initial status invariant
  - **Property 9: Initial Status Invariant**
  - **Validates: Requirements 5.1**

## 3. Checkpoint - Ensure Core Logic Tests Pass

- [ ] 3. Checkpoint



  - Ensure all tests pass, ask the user if questions arise.

## 4. Purchase Order Status Management

- [ ] 4.1 Add status transition validation to UpdateOrder
  - Implement valid status transition rules
  - Add error handling for invalid transitions
  - _Requirements: 5.2, 5.3, 5.4_

- [ ] 4.2 Enhance DeleteOrder with status check
  - Verify current implementation correctly rejects non-pending/cancelled orders
  - Add detailed error message for invalid delete attempts
  - _Requirements: 5.5_

- [ ] 4.3 Write property test for delete constraint
  - **Property 10: Delete Constraint**
  - **Validates: Requirements 5.5**

## 5. Purchase Order Query Enhancement

- [ ] 5.1 Enhance GetOrdersBySupplier response structure
  - Add supplier name to grouped response
  - Include product details (name, unit) in each item
  - _Requirements: 4.1, 4.2, 4.3_

- [ ] 5.2 Write property test for supplier grouping correctness
  - **Property 8: Supplier Grouping Correctness**
  - **Validates: Requirements 4.1, 4.2**

- [ ] 5.3 Verify list query filters work correctly
  - Test store_id, supplier_id, status, date range filters
  - Ensure pagination works as expected
  - _Requirements: 6.1, 6.2, 6.3_

- [ ] 5.4 Write property test for filter correctness
  - **Property 11: Filter Correctness**
  - **Validates: Requirements 6.1**

- [ ] 5.5 Write property test for pagination correctness
  - **Property 12: Pagination Correctness**
  - **Validates: Requirements 6.2**

## 6. Checkpoint - Ensure All Tests Pass

- [ ] 6. Checkpoint
  - Ensure all tests pass, ask the user if questions arise.

## 7. API Routes and Controller Updates

- [ ] 7.1 Add new API routes for store supplier management
  - POST `/api/store-suppliers/bind` - bind products
  - DELETE `/api/store-suppliers/unbind` - unbind products
  - PUT `/api/store-suppliers/default` - set default supplier
  - GET `/api/store-suppliers/:store_id` - list store bindings
  - _Requirements: 1.1, 1.3, 2.1_

- [ ] 7.2 Add API route for supplier-grouped order view
  - GET `/api/purchase-orders/:id/by-supplier` - get order grouped by supplier
  - _Requirements: 4.1_

- [ ] 7.3 Write unit tests for controller endpoints
  - Test request validation
  - Test error responses
  - Test successful responses
  - _Requirements: 1.1, 3.1, 4.1_

## 8. Final Checkpoint - Ensure All Tests Pass

- [ ] 8. Final Checkpoint
  - Ensure all tests pass, ask the user if questions arise.
