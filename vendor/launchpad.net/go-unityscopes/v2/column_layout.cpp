#include <unity/scopes/ColumnLayout.h>
#include <unity/scopes/Variant.h>

#include <unity/UnityExceptions.h>

extern "C" {
#include "_cgo_export.h"
}

#include <cstring>

#include "helpers.h"
#include "smartptr_helper.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

void destroy_column_layout(_ColumnLayout *layout) {
    delete reinterpret_cast<ColumnLayout*>(layout);
}

_ColumnLayout *new_column_layout(int num_columns) {
    return reinterpret_cast<_ColumnLayout*>(new ColumnLayout(num_columns));
}

void column_layout_add_column(_ColumnLayout *layout, const StrData widget_list, char **error) {
    std::vector<std::string> widgets = split_strings(widget_list);
    try {
        reinterpret_cast<ColumnLayout*>(layout)->add_column(widgets);
    } catch(unity::LogicException & e) {
        *error = strdup(e.what());
    }
}

int column_layout_number_of_columns(_ColumnLayout *layout) {
    return reinterpret_cast<ColumnLayout*>(layout)->number_of_columns();
}

int column_layout_size(_ColumnLayout *layout) {
    return reinterpret_cast<ColumnLayout*>(layout)->size();
}

void *column_layout_column(_ColumnLayout *layout, int column, int *length, char **error) {
    try {
        auto columns = reinterpret_cast<ColumnLayout*>(layout)->column(column);
        VariantArray var_array;
        for (auto item: columns) {
            var_array.push_back(Variant(item));
        }
        std::string json_data = Variant(var_array).serialize_json();
        return as_bytes(json_data, length);

    } catch(const std::exception & e) {
        *error = strdup(e.what());
        return nullptr;
    }
}
