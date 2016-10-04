#include <cstring>

#include <unity/scopes/ChildScope.h>
#include <unity/scopes/ScopeMetadata.h>

extern "C" {
#include "_cgo_export.h"
}
#include "helpers.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

_ChildScope *new_child_scope(const StrData id, _ScopeMetadata *metadata, int enabled, const StrData keyword_list) {
    ScopeMetadata *api_metadata = reinterpret_cast<ScopeMetadata *>(metadata);

    std::set<std::string> keywords;
    for (auto &k : split_strings(keyword_list)) {
        keywords.emplace(std::move(k));
    }

    return reinterpret_cast<_ChildScope *>(new ChildScope(from_gostring(id), *api_metadata, enabled, keywords));
}

void destroy_child_scope(_ChildScope *childscope) {
    delete reinterpret_cast<ChildScope*>(childscope);
}

char *child_scope_get_id(_ChildScope *childscope) {
    return strdup(reinterpret_cast<ChildScope*>(childscope)->id.c_str());
}

void set_child_scopes_list(void *child_scopes_list, _ChildScope **source_child_scopes, int length) {
    ChildScopeList *c_child_scopes_list = reinterpret_cast<ChildScopeList*>(child_scopes_list);
    for (int i=0; i < length; ++i) {
        ChildScope *pItem = reinterpret_cast<ChildScope*>(source_child_scopes[i]);
        c_child_scopes_list->push_back(*pItem);
    }
}
