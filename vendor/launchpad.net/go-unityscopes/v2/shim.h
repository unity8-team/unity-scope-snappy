#ifndef UNITYSCOPE_SHIM_H
#define UNITYSCOPE_SHIM_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/* A typedef that can be used to represent a std::shared_ptr */

typedef uintptr_t SharedPtrData[2];

typedef struct StrData {
    char *data;
    long length;
} StrData;

typedef struct _CannedQuery _CannedQuery;
typedef struct _Result _Result;
typedef struct _Result _CategorisedResult;
typedef struct _SearchMetadata _SearchMetadata;
typedef struct _ActionMetadata _ActionMetadata;
typedef struct _ScopeMetadata _ScopeMetadata;
typedef struct _QueryMetadata _QueryMetadata;
typedef struct _ColumnLayout _ColumnLayout;
typedef struct _ChildScope _ChildScope;
typedef void _ScopeBase;
typedef struct _GoString _GoString;

typedef struct _ActivationResponse _ActivationResponse;

void run_scope(const StrData scope_name, const StrData runtime_config,
               const StrData scope_config, void *pointer_to_iface,
               char **error);

/* ScopeBase objects */
char *scope_base_scope_directory(_ScopeBase *scope);
char *scope_base_cache_directory(_ScopeBase *scope);
char *scope_base_tmp_directory(_ScopeBase *scope);
void *scope_base_settings(_ScopeBase *scope, int *length);
_ScopeMetadata **list_registry_scopes_metadata(_ScopeBase *scope, int *n_scopes);
_ChildScope **list_child_scopes(_ScopeBase *scope, int *n_scopes);

/* ChildScope objects */
_ChildScope *new_child_scope(const StrData id, _ScopeMetadata *metadata, int enabled, const StrData keyword_list);
void destroy_child_scope(_ChildScope *childscope);
char *child_scope_get_id(_ChildScope *childscope);
void set_child_scopes_list(void *child_scopes_list, _ChildScope **source_child_scopes, int length);

/* SearchReply objects */
void init_search_reply_ptr(SharedPtrData dest, SharedPtrData src);
void destroy_search_reply_ptr(SharedPtrData data);

void search_reply_finished(SharedPtrData reply);
void search_reply_error(SharedPtrData reply, const StrData err_string);
void search_reply_register_category(SharedPtrData reply, const StrData id, const StrData title, const StrData icon, const StrData cat_template, SharedPtrData category);
void search_reply_register_departments(SharedPtrData reply, SharedPtrData dept);
void search_reply_push(SharedPtrData reply, _CategorisedResult *result, char **error);
void search_reply_push_filters(SharedPtrData reply, const StrData filters_json, const StrData filter_state_json, char **error);

/* PreviewReply objects */
void init_preview_reply_ptr(SharedPtrData dest, SharedPtrData src);
void destroy_preview_reply_ptr(SharedPtrData data);

void preview_reply_finished(SharedPtrData reply);
void preview_reply_error(SharedPtrData reply, const StrData err_string);
void preview_reply_push_widgets(SharedPtrData reply, const StrData widget_list, char **error);
void preview_reply_push_attr(SharedPtrData reply, const StrData key, const StrData json_value, char **error);
void preview_reply_register_layout(SharedPtrData reply, _ColumnLayout **layout, int n_items, char **error);

/* CannedQuery objects */
void destroy_canned_query(_CannedQuery *query);
_CannedQuery *new_canned_query(const StrData scope_id, const StrData query_str, const StrData department_id);
char *canned_query_get_scope_id(_CannedQuery *query);
char *canned_query_get_department_id(_CannedQuery *query);
char *canned_query_get_query_string(_CannedQuery *query);
void *canned_query_get_filter_state(_CannedQuery *query, int *length);
void canned_query_set_department_id(_CannedQuery *query, const StrData department_id);
void canned_query_set_query_string(_CannedQuery *query, const StrData query_str);
char *canned_query_to_uri(_CannedQuery *query);

/* Category objects */
void destroy_category_ptr(SharedPtrData data);

/* CategorisedResult objects */
_Result *new_categorised_result(SharedPtrData category);
void destroy_result(_Result *res);

/* Result objects */
void *result_get_attr(_Result *res, const StrData attr, int *length, char **error);
void result_set_attr(_Result *res, const StrData attr, const StrData json_value, char **error);
void result_set_intercept_activation(_Result *res);

/* Department objects */
void init_department_ptr(SharedPtrData dest, SharedPtrData src);
void new_department(const StrData deptt_id, _CannedQuery *query, const StrData label, SharedPtrData dept, char **error);
void destroy_department_ptr(SharedPtrData data);
void department_add_subdepartment(SharedPtrData dept, SharedPtrData child);
void department_set_alternate_label(SharedPtrData dept, const StrData label);
char *department_get_alternate_label(SharedPtrData dept);
char *department_get_id(SharedPtrData dept);
char *department_get_label(SharedPtrData dept);
void department_set_has_subdepartments(SharedPtrData dept, int subdepartments);
int department_has_subdepartments(SharedPtrData dept);
//void department_get_subdepartments(SharedPtrData dept, SharedPtrData **ret_data);
SharedPtrData * department_get_subdepartments(SharedPtrData dept, int *n_subdepts);
void department_set_subdepartments(SharedPtrData dept, SharedPtrData *subdepartments, int nb_subdepartments);
_CannedQuery * department_get_query(SharedPtrData dept);

/* QueryMetadata objects */
char *query_metadata_get_locale(_QueryMetadata *metadata);
char *query_metadata_get_form_factor(_QueryMetadata *metadata);
void query_metadata_set_internet_connectivity(_QueryMetadata *metadata, int status);
int query_metadata_get_internet_connectivity(_QueryMetadata *metadata);

/* SearchMetadata objects */
_SearchMetadata *new_search_metadata(int cardinality, const StrData locale, const StrData form_factor);
void destroy_search_metadata(_SearchMetadata *metadata);
int search_metadata_get_cardinality(_SearchMetadata *metadata);
void *search_metadata_get_location(_SearchMetadata *metadata, int *length);
void search_metadata_set_location(_SearchMetadata *metadata, char *json_data, int json_data_length, char **error);
void search_metadata_set_aggregated_keywords(_SearchMetadata *metadata, const StrData keyword_list, char **error);
void *search_metadata_get_aggregated_keywords(_SearchMetadata *metadata, int *length);
int search_metadata_is_aggregated(_SearchMetadata *metadata);

/* ActionMetadata objects */
_ActionMetadata *new_action_metadata(const StrData locale, const StrData form_factor);
void destroy_action_metadata(_ActionMetadata *metadata);
void *action_metadata_get_scope_data(_ActionMetadata *metadata, int *data_length);
void action_metadata_set_scope_data(_ActionMetadata *metadata, char *json_data, int json_data_length, char **error);
void action_metadata_set_hint(_ActionMetadata *metadata, const StrData key, char *json_data, int json_data_length, char **error);
void *action_metadata_get_hint(_ActionMetadata *metadata, const StrData key, int *data_length, char **error);
void *action_metadata_get_hints(_ActionMetadata *metadata, int *length);

/* ScopeMetadata objects */
void destroy_scope_metadata_ptr(_ScopeMetadata *metadata);
char *get_scope_metadata_serialized(_ScopeMetadata *metadata);

/* ActivationResponse objects */
void activation_response_init_status(_ActivationResponse *response, int status);
void activation_response_init_query(_ActivationResponse *response, _CannedQuery *query);
void activation_response_init_update_result(_ActivationResponse *response, _Result *result);
void activation_response_init_update_preview(_ActivationResponse *response, const StrData widget_list, char **error);
void activation_response_set_scope_data(_ActivationResponse *response, char *json_data, int json_data_length, char **error);

/* ColumnLayout objects */
_ColumnLayout *new_column_layout(int num_columns);
void destroy_column_layout(_ColumnLayout *layout);
void column_layout_add_column(_ColumnLayout *layout, const StrData widget_list, char **error);
int column_layout_number_of_columns(_ColumnLayout *layout);
int column_layout_size(_ColumnLayout *layout);
void *column_layout_column(_ColumnLayout *layout, int column, int *n_items, char **error);


/* Helpers for tests */
_Result *new_testing_result(void);


#ifdef __cplusplus
}
#endif

#endif
